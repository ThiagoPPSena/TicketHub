# Número de requisições simultâneas que deseja enviar
$numRequests = 20

# URL do endpoint para onde as requisições serão enviadas
$url = "http://localhost:8080/passages/buy"

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
    # Inicia um novo job em background para cada requisição
    $job = Start-Job -ScriptBlock {
        param ($url, $data)
        
        # Envia a requisição HTTP POST
        Invoke-RestMethod -Uri $url -Method Post -Body $data -ContentType "application/json"
    } -ArgumentList $url, $data
    
    # Adiciona o job à lista de jobs
    $jobs += $job
}

# Aguarda que todos os jobs sejam concluídos e exibe os resultados
$jobs | ForEach-Object {
    # Espera o job ser concluído antes de receber o resultado
    Wait-Job -Job $_

    $result = Receive-Job -Job $_
    Write-Output "Resultado da requisição $_.Id: $result"
    
    # Remove o job
    Remove-Job -Job $_
}
