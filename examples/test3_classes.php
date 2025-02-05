<?php
class DataProcessor {
    private $data;
    
    public function __construct(array $data) {
        $this->data = $data;
        usleep(1000);
    }
    
    public function process() {
        $result = 0;
        foreach ($this->data as $item) {
            $result += $this->heavyCalculation($item);
        }
        return $result;
    }
    
    private function heavyCalculation($value) {
        usleep(500);
        return $value * 2;
    }
}

$processor = new DataProcessor(range(1, 5));
echo $processor->process();
