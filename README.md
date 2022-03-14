# Image Editor
[![forthebadge made-with-go](http://ForTheBadge.com/images/badges/made-with-go.svg)](https://go.dev/)
[![GitHub](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/)
[![Discord](https://img.shields.io/badge/Discord-7289DA?style=for-the-badge&logo=discord&logoColor=white)](https://discord.com/)
- Universidad de La Laguna
- **Subject:** Computer vision
- **Final Project of Subject:** Image editor prototype

## Index
- [Authors](#authors)
- [Description](#description)
- [Resources used](#recursos-used)
- [Directory structure](#directory-structure)

## Authors
  - Carlos García Lezcano - alu0101208268@ull.edu.es
  - Eduardo Expósito Barrera - alu0101230382@ull.edu.es

## Description
  - This repository contains a final project of the subject Computer Vision. The project consisted of a prototype of an image editor with different operations seen in the subject.

## Resources used
- The following resources have been used to carry out the project:
    - **[Go](https://go.dev/):** Programming language used for development.
    <br>
    <p align="center">
      <img src="img/Go.png" width="300px">
    </p> 
    
    - **[Github](https://github.com/):** Software used to maintain version control of the developed code.
    <br>
    <p align="center">
      <img src="img/Git.png" width="300px">
    </p>
    
    - **[Discord](https://discord.com/):** Application used for equipment communication.
    <br>
    <p align="center">
      <img src="img/Discord.png" width="200px">
    </p> 


## Directory structure
- The directory is organized as follows:

      .
      ├── img
          ├── Discord.png
          ├── Git.png
          ├── Go.png
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
      ├── testImages
          ├── lena.tiff
          ├── lena2.tiff
          ├── lena3.tiff
          ├── tanque-anterior.tiff
          ├── tanque-posterior.tiff
      ├── go.mod
      ├── go.sum
      ├── main.go