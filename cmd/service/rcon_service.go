package service

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"venecraft-back/cmd/dto"

	"github.com/gorcon/rcon"
)

func SendCommand(client *rcon.Conn, command string) (*dto.CommandResponse, error) {
	// Ejecutar comando RCON
	output, err := client.Execute(command)
	if err != nil {
		return nil, err
	}

	return &dto.CommandResponse{Output: output}, nil
}

type ConsoleLogService struct {
	logFilePath string
}

func NewConsoleLogService() *ConsoleLogService {
	return &ConsoleLogService{
		logFilePath: "http://10.147.17.34:5500/logs/latest.log", // Cambia esta ruta si es necesario
	}
}

// Streaming de logs en tiempo real
func (cls *ConsoleLogService) StreamLogs(w http.ResponseWriter) error {
	// Verifica si la ruta es una URL
	if strings.HasPrefix(cls.logFilePath, "http://") || strings.HasPrefix(cls.logFilePath, "https://") {
		return cls.streamLogsFromURL(w)
	}

	// Si no es una URL, trata de abrir como archivo local
	return cls.streamLogsFromFile(w)
}

// Streaming de logs desde una URL
func (cls *ConsoleLogService) streamLogsFromURL(w http.ResponseWriter) error {
	resp, err := http.Get(cls.logFilePath)
	if err != nil {
		return errors.New("no se pudo obtener los logs desde la URL" + cls.logFilePath + ": " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("error al obtener los logs: " + resp.Status)
	}

	// Configurar la respuesta HTTP para streaming
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	// Leer los logs línea por línea y enviarlos al cliente
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.New("error al leer los logs desde la URL")
		}

		_, writeErr := w.Write([]byte(line))
		if writeErr != nil {
			return errors.New("error al enviar los datos")
		}
		w.(http.Flusher).Flush()
	}

	return nil
}

// Streaming de logs desde un archivo local
func (cls *ConsoleLogService) streamLogsFromFile(w http.ResponseWriter) error {
	file, err := os.Open(cls.logFilePath)
	if err != nil {
		return errors.New("no se pudo abrir el archivo de logs local")
	}
	defer file.Close()

	// Configurar la respuesta HTTP para streaming
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	// Leer los logs línea por línea y enviarlos al cliente
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.New("error al leer los logs desde el archivo local")
		}

		_, writeErr := w.Write([]byte(line))
		if writeErr != nil {
			return errors.New("error al enviar los datos")
		}
		w.(http.Flusher).Flush()
	}

	return nil
}
