<?php
function factorial($n) {
    if ($n <= 1) return 1;
    usleep(500); // Simulate work
    return $n * factorial($n - 1);
}
echo factorial(10);