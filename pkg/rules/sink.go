package rules

import (
	"strings"

	"github.com/nuclio/logger"
	"github.com/valyala/fasthttp"
)

type HTTPSink struct {
	output *HTTPOutput
	log    logger.Logger

	request  *fasthttp.Request
	response *fasthttp.Response
	input    chan []byte
}

func (s *HTTPSink) Start() {
	s.request = fasthttp.AcquireRequest()
	s.response = fasthttp.AcquireResponse()
	s.input = make(chan []byte, 100)

	s.log.DebugWith("Updating request URI", "uri", s.output.Endpoint)
	s.request.URI().Update(s.output.Endpoint)
	s.request.Header.SetContentType("application/json")
	s.log.DebugWith("Updating request Method", "method", s.output.Method)
	s.request.Header.SetMethod(strings.ToUpper(s.output.Method))
	for key, value := range s.output.Headers {
		s.log.DebugWith("Updating request Header", "header", key, "value", value)
		s.request.Header.Add(key, value)
	}

	for key, value := range s.output.Params {
		s.log.DebugWith("Updating query args", "name", key, "value", value)
		s.request.URI().QueryArgs().Add(key, value)
	}

	for key, value := range s.output.Authentication.Header {
		s.log.DebugWith("Updating authentication Header", "header", key, "value", value)
		s.request.Header.Add(key, value)
	}

	s.log.Debug("Starting goroutine to update body")
	go func() {
		for data := range s.input {
			s.request.SetBody(data)
			if err := fasthttp.Do(s.request, s.response); err != nil {
				s.log.WarnWith("sending request error", "error", err)
			}
		}
	}()
}

func (s *HTTPSink) Send(data []byte) {
	s.input <- data
}
