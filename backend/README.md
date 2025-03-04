# Experiview backend

Currently contains two API handlers:

-   GET /experiments?filter={filter}&startDate={startDate}&endDate={endDate} - Gets and filters experiments based on startDate, endDate and a specified filter for the name of the experiment
-   POST /experiments - Request body: `
{
name: string,
date: string,
description: string
}`
