# go-dbcompare

`go-dbcompare` is a Go-based tool designed to compare two PostgreSQL databases and identify differences in table structures and data. The results can be exported to an Excel file for further analysis.

## Features

- Compare tables between two databases.
- Identify missing or extra records in either database.
- Export comparison results to an Excel file.

## Installation

To install and use `go-dbcompare`, follow these steps:

### Clone the Repository
```sh
git clone https://github.com/CDavidSV/go-dbcompare.git
cd go-dbcompare
```

### Build the Project
```sh
go build -o dbcompare main.go
```

## Usage

### Configure Database Connections
Create a configuration file to provide connection details for the two databases.

`db-compare-config.json`
```json
{
    "database1": {
        "name": "Production Database",
        "host": "db1.example.com",
        "port": 5432,
        "database": "postgres",
        "username": "admin",
        "password": "1234"
    },
    "database2": {
        "name": "Development Database",
        "host": "db2.example.com",
        "port": 5432,
        "database": "postgres",
        "username": "admin",
        "password": "1234"
    }
}

```

### Run the Comparison
```sh
./dbcompare compare -o "./results"
```

## Future Improvements
- Schema comparison for detecting index and constraint differences.
- Improved support for multiple database systems.
- Command-line enhancements for better logging and filtering options.
- Additional output formats (HTML, PDF reports).

## Contributing
Contributions are welcome! Feel free to open issues or submit pull requests.

## License
This project is licensed under the MIT License. See the LICENSE file for details.

