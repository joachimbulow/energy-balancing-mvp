#!/bin/bash

echo '
$$$$$$$\  $$\   $$\ $$$$$$\ $$\       $$$$$$$\         $$$$$$\  $$\   $$\ $$$$$$$\        $$$$$$$\  $$\   $$\  $$$$$$\  $$\   $$\ 
$$  __$$\ $$ |  $$ |\_$$  _|$$ |      $$  __$$\       $$  __$$\ $$$\  $$ |$$  __$$\       $$  __$$\ $$ |  $$ |$$  __$$\ $$ |  $$ |
$$ |  $$ |$$ |  $$ |  $$ |  $$ |      $$ |  $$ |      $$ /  $$ |$$$$\ $$ |$$ |  $$ |      $$ |  $$ |$$ |  $$ |$$ /  \__|$$ |  $$ |
$$$$$$$\ |$$ |  $$ |  $$ |  $$ |      $$ |  $$ |      $$$$$$$$ |$$ $$\$$ |$$ |  $$ |      $$$$$$$  |$$ |  $$ |\$$$$$$\  $$$$$$$$ |
$$  __$$\ $$ |  $$ |  $$ |  $$ |      $$ |  $$ |      $$  __$$ |$$ \$$$$ |$$ |  $$ |      $$  ____/ $$ |  $$ | \____$$\ $$  __$$ |
$$ |  $$ |$$ |  $$ |  $$ |  $$ |      $$ |  $$ |      $$ |  $$ |$$ |\$$$ |$$ |  $$ |      $$ |      $$ |  $$ |$$\   $$ |$$ |  $$ |
$$$$$$$  |\$$$$$$  |$$$$$$\ $$$$$$$$\ $$$$$$$  |      $$ |  $$ |$$ | \$$ |$$$$$$$  |      $$ |      \$$$$$$  |\$$$$$$  |$$ |  $$ |
\_______/  \______/ \______|\________|\_______/       \__|  \__|\__|  \__|\_______/       \__|       \______/  \______/ \__|  \__|                                               
                                                                                        
'

# Make sure the script exits if any command fails :d
set -e

# This process is also documented in Notion :)
read -p "Enter the image version tag: " VERSION

# 1. Build the Maven project
echo "Building Maven project..."
mvn clean package

# 2. Build the Docker image
echo "Building Docker image..."
docker build --platform linux/amd64 -t coordinator .

# 3. Tag the Docker image with the specified version
echo "Tagging Docker image with version $VERSION..."
docker tag coordinator pemmvpregistry.azurecr.io/coordinator:v$VERSION

# 4. Log in to the Azure Container Registry using an access token
echo "Logging in to Azure Container Registry..."
TOKEN=$(az acr login --name pemmvpregistry --expose-token --output tsv --query accessToken)
docker login pemmvpregistry.azurecr.io --username 00000000-0000-0000-0000-000000000000 --password $TOKEN

# 5. Push the Docker image to the Azure Container Registry
echo "Pushing Docker image to Azure Container Registry..."
docker push pemmvpregistry.azurecr.io/coordinator:v$VERSION

echo "Done!"
echo "Your image tag name is: pemmvpregistry.azurecr.io/coordinator:v$VERSION"