package crawler

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/cenkalti/backoff/v4"
	"github.com/dmteterin/firefly-assignment/internal/config"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/proxy"
	"github.com/gocolly/colly/queue"
	"github.com/rs/zerolog"
)

type Crawler struct {
	logger           zerolog.Logger
	proxies          []string
	links            []string
	bank             Bank
	queueThreadCount int
	textCh           chan string
	tokenCh          chan string
	errorUrlRetries  map[string]int
	maxRetries       int
	backoffEnabled   bool
	collector        *colly.Collector
}

type Bank interface {
	CountMatches(tokens chan string)
}

func New(cfg *config.Config, bank Bank, logger zerolog.Logger) (*Crawler, error) {
	linksFile, err := os.Open(cfg.CrawlerLinkListPath)
	if err != nil {
		return nil, err
	}

	linksFileScanner := bufio.NewScanner(linksFile)
	linksFileScanner.Split(bufio.ScanLines)

	var links []string
	for linksFileScanner.Scan() {
		links = append(links, linksFileScanner.Text())
	}

	proxiesFile, err := os.Open(cfg.CrawlerProxyListPath)
	if err != nil {
		return nil, err
	}

	proxiesFileScanner := bufio.NewScanner(proxiesFile)
	proxiesFileScanner.Split(bufio.ScanLines)

	var proxies []string
	for proxiesFileScanner.Scan() {
		proxies = append(proxies, cfg.CrawlerProxyTypePrefix+proxiesFileScanner.Text())
	}

	return &Crawler{
		logger:           logger,
		proxies:          proxies,
		links:            links,
		queueThreadCount: cfg.CrawlerQueueConsumerThreadCount,
		textCh:           make(chan string),
		tokenCh:          make(chan string),
		bank:             bank,
		backoffEnabled:   cfg.CrawlerEnableExpBackoff,
		errorUrlRetries:  make(map[string]int),
		maxRetries:       cfg.CrawlerURLRetryLimit,
		collector: colly.NewCollector(
			colly.CacheDir("./cache"),
		),
	}, nil
}

func (c *Crawler) Tokenize() {
	go c.bank.CountMatches(c.tokenCh)

	articleCount := 0

	for text := range c.textCh {
		tokens := strings.FieldsFunc(text, func(r rune) bool {
			return !unicode.IsLetter(r)
		})
		for _, token := range tokens {
			c.tokenCh <- token
		}
		articleCount++
		if articleCount%1000 == 0 {
			c.logger.Info().Msgf("Scanned %v articles", articleCount)
		}
	}
	close(c.tokenCh)
}

func (c *Crawler) configureCollector() {
	extensions.RandomUserAgent(c.collector)

	if p, err := proxy.RoundRobinProxySwitcher(
		c.proxies...,
	); err == nil {
		c.collector.SetProxyFunc(p)
	}

	c.collector.OnHTML(".caas-content-wrapper", func(e *colly.HTMLElement) {
		var sb strings.Builder

		sb.WriteString(e.ChildText("h1"))
		sb.WriteString(" ")
		sb.WriteString(e.ChildText("h2"))
		sb.WriteString(" ")
		sb.WriteString(e.ChildText("p"))

		c.textCh <- sb.String()
	})

	c.collector.OnError(func(r *colly.Response, err error) {

		if err.Error() == "Not Found" {
			return
		}

		if c.backoffEnabled {
			backoff.Retry(func() error {
				return r.Request.Retry()
			}, backoff.NewExponentialBackOff())

			return
		}

		link := r.Request.URL.String()

		if _, ok := c.errorUrlRetries[link]; !ok {
			c.errorUrlRetries[link]++
			r.Request.Retry()
		} else if c.errorUrlRetries[link] <= c.maxRetries {
			c.errorUrlRetries[link]++
			r.Request.Retry()
		} else {
			fmt.Println(c.errorUrlRetries[link], link)
		}
		r.Request.Retry()
	})
}

func (c *Crawler) RunScrapingQueue() error {
	c.configureCollector()

	go c.Tokenize()

	q, err := queue.New(
		c.queueThreadCount,
		&queue.InMemoryQueueStorage{MaxSize: len(c.links)},
	)
	if err != nil {
		c.logger.Error().Err(err).Msg("Could not start scraping queue")
		return err
	}

	for _, link := range c.links {
		q.AddURL(link)
	}

	q.Run(c.collector)
	close(c.textCh)
	return nil
}
