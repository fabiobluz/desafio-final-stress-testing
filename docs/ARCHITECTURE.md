# ARCHITECTURE

- Pool fixo de workers (concorrência)
- Fila de jobs (canal buffered)
- http.Client tunado (keep-alive, idle pool)
- Coleta e renderização do relatório (texto/JSON)
