package auxiliary

import (
	"context"
	"ultima/model"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type Ultima struct {
	UltimaContext context.Context
}

func Init() (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("excludeSwitches", "enable-automation"),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("enable-automation", false),
	)

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)

	ctx, cancel := chromedp.NewContext(allocCtx)
	return ctx, cancel
}

func (ultima *Ultima) ApplySettings() error {
	err := chromedp.Run(ultima.UltimaContext,
		network.Enable(),
		fetch.Enable().WithPatterns([]*fetch.RequestPattern{
			{
				URLPattern:   "*",
				RequestStage: fetch.RequestStageRequest,
			},
		}),
	)
	return err
}

func (ultima *Ultima) ApplyHeader(header []Header) {
	chromedp.ListenTarget(ultima.UltimaContext, func(ev any) {
		switch ev := ev.(type) {
		case *network.EventRequestWillBeSent:

		case *fetch.EventRequestPaused:
			go func() {
				originalHeaders := ConvertHeaders(ev.Request.Headers)

				modifiedHeaders := AddAndOverwriteHeaders(header, originalHeaders)

				chromedp.Run(ultima.UltimaContext,
					fetch.ContinueRequest(ev.RequestID).
						WithHeaders(modifiedHeaders),
				)
				// TODO: error handling
			}()
		}
	})
}

func (ultima *Ultima) Run(url string, e *error) {
	err := chromedp.Run(ultima.UltimaContext,
		chromedp.Navigate(url),
	)
	*e = err
	select {}
}

func StartChrome(startURL, configKeyHash, requestHash string) tea.Cmd {
	return func() tea.Msg {
		var errRun error
		ultx, cancel := Init()
		ultima := Ultima{
			UltimaContext: ultx,
		}
		defer cancel()
		if err := ultima.ApplySettings(); err != nil {
			panic(err)
		}

		customHeaders := []Header{
			{
				Key: "User-Agent",
				Val: "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.6778.205 SEB/3.9.0 (x64)",
			},
			{
				Key: "X-SafeExamBrowser-ConfigKeyHash",
				Val: configKeyHash,
			},
			{
				Key: "X-SafeExamBrowser-RequestHash",
				Val: requestHash,
			},
		}

		ultima.ApplyHeader(customHeaders)
		ultima.Run(startURL, &errRun)
		return model.ChromeHandler{}
	}
}
