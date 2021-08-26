#### Proyecto
La estructura del proyecto está basado en clean code. Para ello defino los siguientes paquetes:

- cmd: donde se encuentra la entrada principal de la app (main.go)
- config: donde se encuentra las configuraciones de sistema y funcionales
- domain: donde se encuentran las entidades y los casos de uso
- adapter: donde se encuentran los adaptadores para la infraestructura que utilizaran los casos de uso mediante interfaces
- logger: esto al ser cross a toda la aplicación lo agregue separado
- infraestructura: donde se encuentran todo los soportes externos a la app (base de datos y migraciones en este caso)
- archivo config.yml donde se encuentran todas las configuraciones (Esta es una decisión personal de cómo llevar a cabo la gestión de las configuraciones, no lo he visto en otros lados pero me gusta
separar lo que es son configuraciones de soporte al system de las reglas propias del negocio; por esto es que dentro del archivo config separo en el
settings y businessrules).

#### Tecnologias y frameworks:
- go 1.15
- gomock
- testify
- go-chi
- zap
- migrate
- postgresql

#### Resolucion:
- Cada uno de los endpoints esta separado en casos de uso.
- El sistema corre alrededor de 100 test (si se considera cada parte del adapter por separado).
- Se utiliza JWT para la authorizacion del usuario.
- El password es guardado en sha256.
- La secret para la firma del token encuentra en el archivo de configuracion junto con la expiracion en minutos. 


#### Mejoras pendientes:
- Por defecto se tiene el sslmode en disable para la db, esto deberia poder setearse por configuracion
- Mejorar y definir mas errores custom dentro de la app.
- Meter swagger para la documentacion de los endpoints
- Mejorar las migrations para que cree automáticamente la base de datos.
- Tracing para seguimiento en logs  
- Test E2E 


#### Cómo se ejecuta:
1. Se debe crear la base en postgresql (lahaus por defecto). El sistema generar todas las estructuras necesarias
2. Compilar el proyecto `go build -o binlahaus ./cmd/main.go`  
3. Cambiar los permisos a ejecucion `chmod +x binlahaus`
4. Ejecutar la aplicacion 

