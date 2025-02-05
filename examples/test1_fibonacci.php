<?php
function fibonacci($n) {
    if ($n <= 1) return $n;
    $prev = 0;
    $curr = 1;
    for ($i = 2; $i <= $n; $i++) {
        $temp = $curr;
        $curr = $prev + $curr;
        $prev = $temp;
        usleep(1000); // Simulate work
    }
    return $curr;
}
echo fibonacci(20);