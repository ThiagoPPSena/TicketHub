# Número de requisições simultâneas que deseja enviar
$numRequests = 20

# URLs dos endpoints para onde as requisições serão enviadas
$urls = @(
    "http://localhost:8080/passages/buy",
    "http://localhost:8081/passages/buy",
    "http://localhost:8082/passages/buy"
)

# Corpo da requisição (para POST ou PUT) em JSON
$data = @{
    routes = @(
        @{
            From = "ARACAJU"
            To = "SALVADOR"
            Seats = 500
            Company = "A"
        },
        @{
            From = "SALVADOR"
            To = "BRASILIA"
            Seats = 250
            Company = "C"
        },
        @{
            From = "BRASILIA"
            To = "GOIANIA"
            Seats = 500
            Company = "B"
        },
        @{
            From = "GOIANIA"
            To = "MANAUS"
            Seats = 200
            Company = "C"
        },
        @{
            From = "MANAUS"
            To = "PORTO VELHO"
            Seats = 300
            Company = "B"
        }
    )
} | ConvertTo-Json

# Array para armazenar os jobs em segundo plano
$jobs = @()

# Loop para iniciar múltiplas requisições em paralelo
for ($i = 1; $i -le $numRequests; $i++) {
    # Seleciona um dos servidores aleatoriamente
    $url = $urls[$i % $urls.Length]

    # Inicia um novo job em background para cada requisição
    $job = Start-Job -ScriptBlock {
        param ($url, $data)
        
        # Envia a requisição HTTP POST
        try {
            Invoke-RestMethod -Uri $url -Method Post -Body $data -ContentType "application/json"
        } catch {
            # Captura e retorna qualquer erro
            "Erro: $_"
        }
    } -ArgumentList $url, $data
    
    # Adiciona o job à lista de jobs
    $jobs += $job
}

# Aguarda que todos os jobs sejam concluídos e exibe os resultados
$jobs | ForEach-Object {
    # Espera o job ser concluído antes de receber o resultado
    Wait-Job -Job $_

    # Obtém o resultado ou o erro do job
    $result = Receive-Job -Job $_
    Write-Output "Resultado da requisição $_.Id para $($urls[($i-1) % $urls.Length]): $result"
    
    # Remove o job
    Remove-Job -Job $_
}
