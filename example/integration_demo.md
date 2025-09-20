# CSV Export Integration Demo

This document demonstrates how to use the new `write_to_csv` skill in combination with the existing `extract_data` skill for complete data collection workflows.

## Workflow Example

1. **Navigate to a webpage**:
   ```json
   {
     "skill": "navigate_to_url",
     "args": {
       "url": "https://example.com/products"
     }
   }
   ```

2. **Extract data from the page**:
   ```json
   {
     "skill": "extract_data",
     "args": {
       "extractors": [
         {
           "name": "product_name",
           "selector": ".product-title",
           "multiple": true
         },
         {
           "name": "price",
           "selector": ".product-price",
           "multiple": true
         },
         {
           "name": "rating",
           "selector": ".product-rating",
           "attribute": "data-rating",
           "multiple": true
         }
       ],
       "format": "json"
     }
   }
   ```

3. **Write the extracted data to CSV**:
   ```json
   {
     "skill": "write_to_csv",
     "args": {
       "data": [
         {"product_name": "Product A", "price": "$29.99", "rating": "4.5"},
         {"product_name": "Product B", "price": "$39.99", "rating": "4.2"},
         {"product_name": "Product C", "price": "$19.99", "rating": "4.8"}
       ],
       "file_path": "/tmp/products.csv",
       "headers": ["product_name", "price", "rating"],
       "include_headers": true
     }
   }
   ```

## Features Supported

- **Custom Headers**: Specify column order and names
- **Append Mode**: Add to existing CSV files without overwriting
- **Flexible Data**: Handles arrays, objects, and primitive values
- **Error Handling**: Validates data format and file operations
- **Directory Creation**: Automatically creates parent directories

## Use Cases

- **E-commerce Data Collection**: Extract product information, prices, and reviews
- **News Aggregation**: Collect headlines, dates, and article links
- **Financial Data**: Gather stock prices, market data, and trading volumes
- **Contact Information**: Extract business details from directory sites
- **Event Listings**: Collect event names, dates, venues, and prices