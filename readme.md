# Devops CLI: A tool to do DevOps in style

A Terminal/Web UI to interact with devops tools & services, intended for devops/platform engineers.

## Inspiration
This project is inspired from the following softwares terraform, K9s, Lens. 
- Devops CLI aims to bring the agility (that extensively comes from using keyboard) & speed offered by K9s.
- Devops CLI aims to bring the extensibility provided by terraform (in form of plugins)
- Not everything can be done on TUI, Something are better show in a web app like Lens. The above 2 things will also be package in a Web UI for better experience.

## This project aims to achieve the following:
- Improve debugging
- Improve development agility

## [Watch Web App Demo 1](https://youtu.be/DtQnuSDmodg)
## [Watch Web App Demo 2 (new ux)](https://youtu.be/0nEwPfzeikQ)

## Installation
**Linux & Mac**

`curl https://storage.googleapis.com/devops-cli-artifacts/releases/devops/0.5.3/install.sh | bash`

## Usage

**Run Server**

`devops`

**Run Client (Web App)**

![image](https://user-images.githubusercontent.com/24411676/230721653-a57f0eea-7629-4839-ba32-1eb6cb77415f.png)


`On browser go to: http://localhost:9753`

**Run Client (TUI)**

`devops tui`

## Read [Wiki](https://github.com/sharadregoti/devops-cli/wiki) for detailed documentation 

## Supported Plugins
- Kubernetes
- Helm
- Gitlab (WIP)

### Kubernetes Features
- View & search all cluster resources
- Create, Read, Update, Delete any resource
- Describe any resource
- View logs, get shell access of pod

### Helm Features
- View releases & charts
- Perform the following operation on releases: rollback, uninstall, view values, view history, view manifies etc...
