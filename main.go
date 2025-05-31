package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Resultado representa o resultado de uma requisição
type Resultado struct {
	URL        string
	StatusCode int
	Erro       error
}

// verificarURL faz uma requisição HTTP e envia o resultado para o canal
func verificarURL(ctx context.Context, url string, resultadoChan chan<- Resultado) {
	client := &http.Client{
		Timeout: 5 * time.Second, // Timeout para evitar requisições lentas
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		resultadoChan <- Resultado{URL: url, StatusCode: 0, Erro: err}
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		resultadoChan <- Resultado{URL: url, StatusCode: 0, Erro: err}
		return
	}
	defer resp.Body.Close()

	resultadoChan <- Resultado{URL: url, StatusCode: resp.StatusCode, Erro: nil}
}

func main() {
	// Lista de URLs para verificar
	urls := []string{
		"https://www.google.com",
		"https://www.github.com",
		"https://www.instagram.com",
		"https://www.facebook.com", // URL inválida para simular erro
	}

	// Cria um contexto com timeout global de 10 segundos
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Canal para coletar resultados
	resultadoChan := make(chan Resultado, len(urls))

	// WaitGroup para esperar todas as goroutines
	var wg sync.WaitGroup

	// Fan-out: Lança uma goroutine para cada URL
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			verificarURL(ctx, u, resultadoChan)
		}(url)
	}

	// Fan-in: Coleta resultados em uma goroutine separada
	go func() {
		wg.Wait()
		close(resultadoChan) // Fecha o canal após todas as goroutines terminarem
	}()

	// Processa os resultados
	for resultado := range resultadoChan {
		if resultado.Erro != nil {
			fmt.Printf("Erro ao verificar %s: %v\n", resultado.URL, resultado.Erro)
		} else {
			fmt.Printf("URL: %s, Status: %d\n", resultado.URL, resultado.StatusCode)
		}
	}

	fmt.Println("Verificação concluída!")
}