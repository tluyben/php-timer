<?php
function processValue($value) {
    usleep(200);
    
    if ($value < 0) {
        usleep(300);
        return "negative";
    } elseif ($value == 0) {
        usleep(400);
        return "zero";
    } else {
        if ($value > 100) {
            usleep(500);
            return "large";
        } else {
            usleep(600);
            return "small";
        }
    }
}

$values = [-5, 0, 50, 150];
foreach ($values as $value) {
    echo processValue($value) . "\n";
}
