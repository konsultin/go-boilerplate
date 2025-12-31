# timek - Time Library Extension

A production-ready time library wrapper around `carbon` with built-in Indonesian/English support and common formatting utilities.

## Features

- **Multi-language**: Built-in support for Indonesian (ID) and English (EN).
- **Custom Formats**: Comprehensive set of pre-defined constants for date/time layouts.
- **Human Readable**: `ToReadable()` for relative time strings (e.g., "2 jam yang lalu").
- **Database/JSON**: Fully compatible with `database/sql` and `encoding/json`.

## Quick Start

```go
import "github.com/konsultin/project-goes-here/libs/timek"

// Get current time (Defaults to ID locale)
now := timek.Now()

// Format as "DD/MM/YY" (e.g., 15/12/24)
fmt.Println(now.ToDate())

// Format as "DD MMMM YYYY" (e.g., 15 Desember 2024)
fmt.Println(now.ToDateFull())

// Humans readable (e.g., "beberapa detik yang lalu")
fmt.Println(now.ToReadable())
```

## Formats Reference

All formats are available as constants. You can use `timek.Format(timek.ConstantName)` or helper methods.

| Constant | Format Code | Output Example |
|----------|-------------|----------------|
| **Date** | | |
| `FormatDate` | `d/m/y` | 15/12/24 |
| `FormatDateStructure` | `Y-m-d` | 2024-12-15 |
| `FormatDateUnStructure`| `d-m-Y` | 15-12-2024 |
| `FormatDateFull` | `d F Y` | 15 Desember 2024 |
| `FormatDateDay` | `l, d F Y` | Minggu, 15 Desember 2024 |
| `FormatDateShortMonth` | `d M Y` | 15 Des 2024 |
| `FormatDateMonthYear` | `F Y` | Desember 2024 |
| **DateTime** | | |
| `FormatDateTime` | `d/m/y H:i` | 15/12/24 14:30 |
| `FormatDateTimeStandard`| `Y-m-d H:i:s`| 2024-12-15 14:30:00 |
| `FormatDateTimeStructure`| `Y-m-d H:i` | 2024-12-15 14:30 |
| `FormatDateTimeFull` | `d F Y, H:i` | 15 Desember 2024, 14:30 |
| `FormatDateTimeDay` | `l, d F Y H:i`| Minggu, 15 Desember 2024 14:30 |
| `FormatDateTimeISO` | `Y-m-d\TH:i:s.u\Z` | 2024-12-15T14:30:00.000Z |
| **Time** | | |
| `FormatTime` | `H:i` | 14:30 |
| `FormatTimeFull` | `H:i:s` | 14:30:59 |

## Database Usage

`timek.Time` implements `Scanner` and `Valuer` interfaces.

```go
type User struct {
    ID        int        `db:"id"`
    CreatedAt timek.Time `db:"created_at"`
    UpdatedAt timek.Time `db:"updated_at"`
}
```
