package serp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oxylabs/oxylabs-sdk-go/internal"
	"github.com/oxylabs/oxylabs-sdk-go/oxylabs"
)

// ScrapeYandexSearch scrapes yandex with async polling runtime via Oxylabs SERP API
// and yandex_search as source.
func (c *SerpClientAsync) ScrapeYandexSearch(
	query string,
	opts ...*YandexSearchOpts,
) (chan *Resp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), internal.DefaultTimeout)
	defer cancel()

	return c.ScrapeYandexSearchCtx(ctx, query, opts...)
}

// ScrapeYandexSearchCtx scrapes yandex with async polling runtime via Oxylabs SERP API
// and yandex_search as source.
// The provided context allows customization of the HTTP req, including setting timeouts.
func (c *SerpClientAsync) ScrapeYandexSearchCtx(
	ctx context.Context,
	query string,
	opts ...*YandexSearchOpts,
) (chan *Resp, error) {
	httpRespChan := make(chan *http.Response)
	respChan := make(chan *Resp)
	errChan := make(chan error)

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

	// Check the validity of the parameters.
	err := opt.checkParameterValidity()
	if err != nil {
		return nil, err
	}

	// Prepare the payload.
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
		payload["parsing_instructions"] = &opt.ParseInstructions
		customParserFlag = true
	}

	// Marshal.
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		payload["parse"] = true
		return nil, fmt.Errorf("error marshalling payload: %v", err)
	}

	// Get the job ID.
	jobID, err := c.C.GetJobID(jsonPayload)
	if err != nil {
		return nil, err
	}

	// Poll job status.
	go c.C.PollJobStatus(
		ctx,
		jobID,
		opt.PollInterval,
		httpRespChan,
		errChan,
	)

	// Handle error.
	err = <-errChan
	if err != nil {
		return nil, err
	}

	// Unmarshal the http Response and get the response.
	httpResp := <-httpRespChan
	resp, err := GetResp(httpResp, customParserFlag, customParserFlag)
	if err != nil {
		return nil, err
	}

	// Retrieve internal resp and forward it to the
	// resp channel.
	go func() {
		respChan <- resp
	}()

	return respChan, nil
}

// ScrapeYandexUrl scrapes yandex with async polling runtime via Oxylabs SERP API
// and yandex as source.
func (c *SerpClientAsync) ScrapeYandexUrl(
	url string,
	opts ...*YandexUrlOpts,
) (chan *Resp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), internal.DefaultTimeout)
	defer cancel()

	return c.ScrapeYandexUrlCtx(ctx, url, opts...)
}

// ScrapeYandexUrlCtx scrapes yandex with async polling runtime via Oxylabs SERP API
// and yandex as source.
// The provided context allows customization of the HTTP req, including setting timeouts.
func (c *SerpClientAsync) ScrapeYandexUrlCtx(
	ctx context.Context,
	url string,
	opts ...*YandexUrlOpts,
) (chan *Resp, error) {
	httpRespChan := make(chan *http.Response)
	respChan := make(chan *Resp)
	errChan := make(chan error)

	// Check the validity of the URL.
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

	// Check the validity of parameters.
	err = opt.checkParameterValidity()
	if err != nil {
		return nil, err
	}

	// Prepare the payload.
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

	// Get the job ID.
	jobID, err := c.C.GetJobID(jsonPayload)
	if err != nil {
		return nil, err
	}

	// Poll job status.
	go c.C.PollJobStatus(
		ctx,
		jobID,
		opt.PollInterval,
		httpRespChan,
		errChan,
	)

	// Handle error.
	err = <-errChan
	if err != nil {
		return nil, err
	}

	// Unmarshal the http Response and get the response.
	httpResp := <-httpRespChan
	resp, err := GetResp(httpResp, customParserFlag, customParserFlag)
	if err != nil {
		return nil, err
	}

	// Retrieve internal resp and forward it to the
	// resp channel.
	go func() {
		respChan <- resp
	}()

	return respChan, nil
}
