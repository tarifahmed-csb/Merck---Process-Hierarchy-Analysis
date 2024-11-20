# Merck Process Hierarchy Modeling with AWS DynamoDB (In Progress)

## Overview

Engineers in our global manufacturing network extract information from billions of rows of data stored in our manufacturing data lake. This often results in significant time wasted in retrieving the correct information for monitoring and troubleshooting issues. To enhance this process, raw data in the data lake is labeled with process context. This context is currently modeled using an SQL database, and we aim to explore how AWS DynamoDB can improve the performance of our hierarchy modeling and analytics efforts.

## Project Description

The **Merck Process Hierarchy Modeling** project is a command-line application that allows users to interact with a DynamoDB database. The application provides several options for managing and querying data, including populating the database, deleting all items, retrieving all stages, and fetching operations by stage.

### Features

- **Populate Database**: Load initial data into the DynamoDB table.
- **Delete All Items**: Clear all entries from the DynamoDB table.
- **Get All Stages**: Retrieve a list of all stages from the database.
- **Get Operations by Stage**: Fetch operations related to a specific stage in the process.

## Requirements

- Go programming language (version x.x.x or higher)
- AWS SDK for Go (for DynamoDB operations)

## Usage

To compile and run the application, use the following commands in your terminal:

```bash
go build -o merckTable
./merckTable



