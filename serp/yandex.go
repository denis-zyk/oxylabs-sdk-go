package serp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mslmio/oxylabs-sdk-go/internal"
	"github.com/mslmio/oxylabs-sdk-go/oxylabs"
)

// Accepted parameters for yandex.
var YandexSearchAcceptedDomainParameters = []oxylabs.Domain{
	oxylabs.DOMAIN_COM,
	oxylabs.DOMAIN_RU,
	oxylabs.DOMAIN_UA,
	oxylabs.DOMAIN_BY,
	oxylabs.DOMAIN_KZ,
	oxylabs.DOMAIN_TR,
}
var YandexSearchAcceptedLocaleParameters = []oxylabs.Locale{
	oxylabs.LOCALE_EN,
	oxylabs.LOCALE_RU,
	oxylabs.LOCALE_BY,
	oxylabs.LOCALE_DE,
	oxylabs.LOCALE_FR,
	oxylabs.LOCALE_ID,
	oxylabs.LOCALE_KK,
	oxylabs.LOCALE_TT,
	oxylabs.LOCALE_TR,
	oxylabs.LOCALE_UK,
}

// checkParameterValidity checks validity of ScrapeYandexSearch parameters.
func (opt *YandexSearchOpts) checkParameterValidity() error {
	if !internal.InList(opt.Domain, YandexSearchAcceptedDomainParameters) {
		return fmt.Errorf("invalid domain parameter: %s", opt.Domain)
	}

	if opt.Locale != "" && !internal.InList(opt.Locale, YandexSearchAcceptedLocaleParameters) {
		return fmt.Errorf("invalid locale parameter: %s", opt.Locale)
	}

	if !oxylabs.IsUserAgentValid(opt.UserAgent) {
		return fmt.Errorf("invalid user agent parameter: %v", opt.UserAgent)
	}

	if opt.Limit <= 0 || opt.Pages <= 0 || opt.StartPage <= 0 {
		return fmt.Errorf("limit, pages and start_page parameters must be greater than 0")
	}

	if opt.ParseInstructions != nil {
		if err := oxylabs.ValidateParseInstructions(opt.ParseInstructions); err != nil {
			return fmt.Errorf("invalid parse instructions: %w", err)
		}
	}

	return nil
}

// checkParameterValidity checks validity of ScrapeYandexUrl parameters.
func (opt *YandexUrlOpts) checkParameterValidity() error {
	if !oxylabs.IsUserAgentValid(opt.UserAgent) {
		return fmt.Errorf("invalid user agent parameter: %v", opt.UserAgent)
	}

	if opt.Render != "" && !oxylabs.IsRenderValid(opt.Render) {
		return fmt.Errorf("invalid render parameter: %v", opt.Render)
	}

	if opt.ParseInstructions != nil {
		if err := oxylabs.ValidateParseInstructions(opt.ParseInstructions); err != nil {
			return fmt.Errorf("invalid parse instructions: %w", err)
		}
	}

	return nil
}

// YandexSearchOpts contains all the query parameters available for yandex_search.
type YandexSearchOpts struct {
	Domain            oxylabs.Domain
	StartPage         int
	Pages             int
	Limit             int
	Locale            oxylabs.Locale
	GeoLocation       string
	UserAgent         oxylabs.UserAgent
	CallbackUrl       string
	ParseInstructions *map[string]interface{}
	PollInterval      time.Duration
}

// ScrapeYandexSearch scrapes yandex via Oxylabs SERP API with yandex_search as source.
func (c *SerpClient) ScrapeYandexSearch(
	query string,
	opts ...*YandexSearchOpts,
) (*Resp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), internal.DefaultTimeout)
	defer cancel()

	return c.ScrapeYandexSearchCtx(ctx, query, opts...)
}

// ScrapeYandexSearchCtx scrapes yandex via Oxylabs SERP API with yandex_search as source.
// The provided context allows customization of the HTTP req, including setting timeouts.
func (c *SerpClient) ScrapeYandexSearchCtx(
	ctx context.Context,
	query string,
	opts ...*YandexSearchOpts,
) (*Resp, error) {
	// Prepare options.
	opt := &YandexSearchOpts{}
	if len(opts) > 0 && opts[len(opts)-1] != nil {
		opt = opts[len(opts)-1]
	}

	// Set defaults.
	internal.SetDefaultDomain(&opt.Domain)
	internal.SetDefaultStartPage(&opt.StartPage)
	internal.SetDefaultLimit(&opt.Limit, internal.DefaultLimit_SERP)
	internal.SetDefaultPages(&opt.Pages)
	internal.SetDefaultUserAgent(&opt.UserAgent)

	// Check validity of parameters.
	err := opt.checkParameterValidity()
	if err != nil {
		return nil, err
	}

	// Prepare payload.
	payload := map[string]interface{}{
		"source":          oxylabs.YandexSearch,
		"domain":          opt.Domain,
		"query":           query,
		"start_page":      opt.StartPage,
		"pages":           opt.Pages,
		"limit":           opt.Limit,
		"locale":          opt.Locale,
		"geo_location":    opt.GeoLocation,
		"user_agent_type": opt.UserAgent,
		"callback_url":    opt.CallbackUrl,
	}

	// Add custom parsing instructions to the payload if provided.
	customParserFlag := false
	if opt.ParseInstructions != nil {
		payload["parse"] = true
		payload["parsing_instructions"] = &opt.ParseInstructions
		customParserFlag = true
	}

	// Marshal.
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %v", err)
	}

	// Req.
	httpResp, err := c.C.Req(ctx, jsonPayload, "POST")
	if err != nil {
		return nil, err
	}

	// Unmarshal the http Response and get the response.
	resp, err := GetResp(httpResp, customParserFlag, customParserFlag)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// YandexUrlOpts contains all the query parameters available for yandex.
type YandexUrlOpts struct {
	UserAgent         oxylabs.UserAgent
	Render            oxylabs.Render
	CallbackUrl       string
	ParseInstructions *map[string]interface{}
	PollInterval      time.Duration
}

// ScrapeYandexUrl scrapes a yandex url via Oxylabs SERP API with yandex as source.
func (c *SerpClient) ScrapeYandexUrl(
	url string,
	opts ...*YandexUrlOpts,
) (*Resp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), internal.DefaultTimeout)
	defer cancel()

	return c.ScrapeYandexUrlCtx(ctx, url, opts...)
}

// ScapeYandexUrlCtx scrapes a yandex url via Oxylabs SERP API with yandex as source.
// The provided context allows customization of the HTTP req, including setting timeouts.
func (c *SerpClient) ScrapeYandexUrlCtx(
	ctx context.Context,
	url string,
	opts ...*YandexUrlOpts,
) (*Resp, error) {
	// Check validity of URL.
	err := internal.ValidateUrl(url, "yandex")
	if err != nil {
		return nil, err
	}

	// Prepare options.
	opt := &YandexUrlOpts{}
	if len(opts) > 0 && opts[len(opts)-1] != nil {
		opt = opts[len(opts)-1]
	}

	// Set defaults.
	internal.SetDefaultUserAgent(&opt.UserAgent)

	// Check validity of parameters.
	err = opt.checkParameterValidity()
	if err != nil {
		return nil, err
	}

	// Prepare payload.
	payload := map[string]interface{}{
		"source":          oxylabs.YandexUrl,
		"url":             url,
		"user_agent_type": opt.UserAgent,
		"render":          opt.Render,
		"callback_url":    opt.CallbackUrl,
	}

	// Add custom parsing instructions to the payload if provided.
	customParserFlag := false
	if opt.ParseInstructions != nil {
		payload["parse"] = true
		payload["parsing_instructions"] = &opt.ParseInstructions
		customParserFlag = true
	}

	// Marshal.
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %v", err)
	}

	// Req.
	httpResp, err := c.C.Req(ctx, jsonPayload, "POST")
	if err != nil {
		return nil, err
	}

	// Unmarshal the http Response and get the response.
	resp, err := GetResp(httpResp, customParserFlag, customParserFlag)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
