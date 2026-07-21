param(
    [switch]$Force
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

$racine = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$fichierEnv = Join-Path $racine '.env'
$dossierCertificats = Join-Path $racine 'traefik\certificats'
$fichierCle = Join-Path $dossierCertificats 'historykanban.key'
$fichierCertificat = Join-Path $dossierCertificats 'historykanban.crt'

$fichiersExistants = @($fichierEnv, $fichierCle, $fichierCertificat) | Where-Object {
    Test-Path -LiteralPath $_
}
if ($fichiersExistants.Count -gt 0 -and -not $Force) {
    throw 'Des secrets existent deja. Relancez avec -Force pour les renouveler.'
}

function New-SecretHexadecimal {
    param([int]$NombreOctets)

    $octets = New-Object byte[] $NombreOctets
    $generateur = [System.Security.Cryptography.RandomNumberGenerator]::Create()
    try {
        $generateur.GetBytes($octets)
    }
    finally {
        $generateur.Dispose()
    }

    return -join ($octets | ForEach-Object { $_.ToString('x2') })
}

$commandeOpenSsl = Get-Command openssl -ErrorAction SilentlyContinue
$candidatsOpenSsl = @(
    if ($commandeOpenSsl) { $commandeOpenSsl.Source }
    'C:\Program Files\Git\usr\bin\openssl.exe'
    'C:\Program Files\Git\mingw64\bin\openssl.exe'
) | Where-Object { $_ -and (Test-Path -LiteralPath $_) }

if ($candidatsOpenSsl.Count -eq 0) {
    throw 'OpenSSL est introuvable. Installez OpenSSL ou Git pour Windows.'
}
$openSsl = $candidatsOpenSsl[0]

New-Item -ItemType Directory -Force -Path $dossierCertificats | Out-Null
$suffixeTemporaire = [Guid]::NewGuid().ToString('N')
$cleTemporaire = "$fichierCle.$suffixeTemporaire.tmp"
$certificatTemporaire = "$fichierCertificat.$suffixeTemporaire.tmp"

try {
    & $openSsl req -x509 -newkey rsa:4096 -sha256 -days 825 -nodes `
        -subj '/CN=historykanban.localhost' `
        -addext 'subjectAltName=DNS:historykanban.localhost,DNS:*.historykanban.localhost,DNS:localhost,IP:127.0.0.1,IP:::1' `
        -keyout $cleTemporaire -out $certificatTemporaire
    if ($LASTEXITCODE -ne 0) {
        throw "OpenSSL a echoue avec le code $LASTEXITCODE."
    }

    $lignesEnv = @(
        "BDMDP=$(New-SecretHexadecimal 24)"
        'BDNOM=historykanban'
        'BDUTILISATEUR=historykanban'
        'KEYCLOAKADMIN=admin'
        "KEYCLOAKADMINMDP=$(New-SecretHexadecimal 24)"
        "MINIOCLE=$(New-SecretHexadecimal 10)"
        "MINIOSECRET=$(New-SecretHexadecimal 20)"
        'RABBITUTILISATEUR=historykanban'
        "RABBITMDP=$(New-SecretHexadecimal 24)"
        'SAUVEGARDEINTERVALLE=86400'
        'SAUVEGARDERETENTION=14'
    )

    [System.IO.File]::WriteAllLines(
        $fichierEnv,
        $lignesEnv,
        [System.Text.UTF8Encoding]::new($false)
    )
    Move-Item -Force -LiteralPath $cleTemporaire -Destination $fichierCle
    Move-Item -Force -LiteralPath $certificatTemporaire -Destination $fichierCertificat
}
finally {
    Remove-Item -Force -LiteralPath $cleTemporaire, $certificatTemporaire -ErrorAction SilentlyContinue
}

Write-Host 'Secrets et certificat locaux renouveles. Aucune valeur sensible n a ete affichee.'
