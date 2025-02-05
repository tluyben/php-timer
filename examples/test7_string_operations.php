<?php
$text = str_repeat("Hello World! ", 1000);

// String operations
$upper = strtoupper($text);
usleep(500);

$lower = strtolower($upper);
usleep(500);

$reversed = strrev($lower);
usleep(500);

$chunks = str_split($reversed, 10);
usleep(500);

echo strlen($text);