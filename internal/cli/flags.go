package cli

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/acme/gload/internal/runner"
)

// Parse processa as flags da CLI e devolve runner.Config pronto para uso.
func Parse(args []string) (runner.Config, error) {
	fs := flag.NewFlagSet("gload", flag.ContinueOnError)
	var (
		urlStr     = fs.String("url", "", "URL alvo (http/https)")
		reqTotal   = fs.Int("requests", 0, "Total de requests")
		conc      = fs.Int("concurrency", 0, "Concorrência (workers)")
		timeout   = fs.Duration("timeout", 10*time.Second, "Timeout por request (ex: 5s, 2m)")
		method    = fs.String("method", "GET", "Método HTTP")
		format    = fs.String("format", "text", "Formato do relatório: text|json")
		body      = fs.String("body", "", "Corpo da requisição (ou @arquivo)")
		headers   multiHeader
	)
	fs.Var(&headers, "header", "Header no formato 'K:V' (pode repetir)")
	if err := fs.Parse(args); err != nil { return runner.Config{}, err }

	if *urlStr == "" { return runner.Config{}, errors.New("--url é obrigatório") }
	u, err := url.Parse(*urlStr)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return runner.Config{}, fmt.Errorf("url inválida: %q", *urlStr)
	}
	if *reqTotal <= 0 { return runner.Config{}, errors.New("--requests deve ser > 0") }
	if *conc <= 0 { return runner.Config{}, errors.New("--concurrency deve ser > 0") }
	if *conc > *reqTotal { *conc = *reqTotal }

	m := strings.ToUpper(*method)
	switch m {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodHead, http.MethodOptions:
	default:
		return runner.Config{}, fmt.Errorf("método não suportado: %s", m)
	}

	var b []byte
	if strings.HasPrefix(*body, "@") {
		path := strings.TrimPrefix(*body, "@")
		data, err := os.ReadFile(path)
		if err != nil { return runner.Config{}, fmt.Errorf("erro lendo body %s: %w", path, err) }
		b = data
	} else {
		b = []byte(*body)
	}

	fmtStr := strings.ToLower(*format)
	if fmtStr != "text" && fmtStr != "json" { return runner.Config{}, fmt.Errorf("--format inválido: %s", *format) }

	return runner.Config{
		URL:         u.String(),
		Requests:    *reqTotal,
		Concurrency: *conc,
		Timeout:     *timeout,
		Method:      m,
		Headers:     headers.Header(),
		Body:        b,
		Format:      fmtStr,
	}, nil
}

type multiHeader struct{ h http.Header }

func (m *multiHeader) String() string { return "" }
func (m *multiHeader) Set(v string) error {
	if m.h == nil { m.h = make(http.Header) }
	parts := strings.SplitN(v, ":", 2)
	if len(parts) != 2 { return fmt.Errorf("header inválido, use K:V") }
	k := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])
	m.h.Add(k, val)
	return nil
}
func (m *multiHeader) Header() http.Header { if m.h == nil { m.h = make(http.Header) }; return m.h }
