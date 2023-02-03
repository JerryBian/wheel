$WWWRootLoc = Join-Path $PSScriptRoot .. static
$NodeModulesLoc = Join-Path $PSScriptRoot .. node_modules

npm install
ncu -u
npm install

Copy-Item -Path $(Join-Path $NodeModulesLoc bootstrap-icons font fonts *) `
    -Destination $(Join-Path $WWWRootLoc fonts)
Write-Output "[all]: Fonts copied"

$jobs = @()

# home page
$jobs += Start-ThreadJob -Name home_css -ScriptBlock {
    npx sass $(Join-Path $using:WWWRootLoc style style.scss) $(Join-Path $using:WWWRootLoc style.css)
    uglifycss --ugly-comments `
        --output $(Join-Path $using:WWWRootLoc style.min.css) `
        $(Join-Path $using:WWWRootLoc style.css)
    Write-Output "CSS completed"
}

$jobs += Start-ThreadJob -Name home_js -ScriptBlock {
    uglifyjs --compress -o $(Join-Path $using:WWWRootLoc script.min.js) `
        $(Join-Path $using:NodeModulesLoc bootstrap dist js bootstrap.bundle.js)
    Write-Output "JS completed"
}

Wait-job -Job $jobs
foreach ($job in $jobs) {
    Receive-Job -Job $job
}

Write-Output "All Done!"