# VPC
[![forthebadge made-with-go](http://ForTheBadge.com/images/badges/made-with-go.svg)](https://go.dev/)
[![GitHub](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/lochdeve/VPC)
- Universidad de La Laguna
- **Asignatura:** Visión por computador
- **Proyecto Final de Asignatura:** Prototipo de editor de imágenes

## Índice
- [Autores](#autores)
- [Descripción](#descripción)
- [Recursos utilizados](#recursos-utilizados)
- [Estructura de directorios](#estructura-de-directorios)

## Autores
  - Carlos García González - alu0101208268@ull.edu.es
  - Eduardo Expósito Barrera - alu0101230382@ull.edu.es

## Descripción
  - Este repositorio contiene un proyecto final de asignatura de la asignatura Visión por computador. El proyecto ha consistido en realizar un prototipo de editor de imáganes con distintas operaciones vistas en la materia.

## Estructura de directorios
- El directorio está organizado de la siguiente manera:

      .
      ├── img
          ├── lena.tiff
          ├── lena2.tiff
          ├── lena3.tiff
          ├── tanque-anterior.tiff
          ├── tanque-posterior.tiff
      ├── pkg
          ├── histogram
              ├── histogram.go
          ├── imageContent
              ├── imageContent.go
          ├── information
              ├── information.go
          ├── loadandsave
              ├── loadandsave.go
          ├── menu
              ├── menu.go
          ├── mouse
              ├── mousEvents.go
          ├── newWindow
              ├── newWindow.go
          ├── operations
              ├── operations.go
      ├── go.mod
      ├── go.sum
      ├── main.go