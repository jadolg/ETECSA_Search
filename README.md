# ETECSA Search
Otro sitio para consultar la base de datos de ETECSA.

### Background
Este sitio fue creado como parte de mi autoestudio utilizando GO en el backend y Vue.js + Materialize en el frontend. No constituye parte de ningún servicio de ETECSA ni se distribuye ninguna información confidencial con el mismo.

### Modo de uso
#### Construir el ejecutable

`go build server.go`

#### Ejecutar el servidor

`./server -db /ruta/a/la/db -port 80`

### Dependencias
 - github.com/labstack/echo
 - github.com/mattn/go-sqlite3
 - github.com/pkg/errors
 - github.com/prometheus/common/log

### Pruebas
Para ejecutar las pruebas navegue hasta el directorio del proyecto y ejecute `go test .`