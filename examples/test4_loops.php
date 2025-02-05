<?php
$array = range(1, 100);
$sum = 0;

// For loop
for ($i = 0; $i < count($array); $i++) {
    $sum += $array[$i];
    usleep(100);
}

// While loop
$i = 0;
while ($i < count($array)) {
    $sum -= $array[$i];
    $i++;
    usleep(100);
}

// Foreach
foreach ($array as $value) {
    $sum += $value * 2;
    usleep(100);
}

echo $sum;