# PHP-Timers 🕒

PHP code profiler that instruments your code to measure execution time. Built using nikic/php-parser.

## Features 🌟

- Automatic code instrumentation
- Original code preservation
- Timing measurements with microsecond precision
- Verbose/aggregated timing modes
- One-command file restoration

## Installation 🛠️

```bash
composer require nikic/php-parser
chmod +x php-timers
```

## Usage 💻

Profile a file:

```bash
./php-timers script.php
```

Detailed profiling:

```bash
./php-timers --verbose script.php
```

Restore originals:

```bash
./php-timers --restore /path/to/directory
```

## How It Works 🔍

1. Creates `filename.__org__.php` backup
2. Instruments code with DateTime measurements
3. Records timing for each statement
4. Outputs results in HTML comments

## Output Format 📋

```php
<!-- Results:
Array
(
    [0] => Array
        (
            [line] => 5
            [code] => "original_code_here"
            [start] => DateTime Object
            [end] => DateTime Object
            [diff] => DateInterval Object
        )
)
-->
```

## Best Practices 🎯

- Use in development only
- --verbose for loop analysis
- Backup important files
- Always restore after profiling

## Limitations ⚠️

- Adds overhead to execution
- Development use only
- Requires file write permissions

## License 📜

MIT License
