<?php
class DB {
    public function query($sql) {
        // Simulate different query times
        if (stripos($sql, 'SELECT') !== false) {
            usleep(1000);
        } else if (stripos($sql, 'INSERT') !== false) {
            usleep(2000);
        } else if (stripos($sql, 'UPDATE') !== false) {
            usleep(3000);
        }
        return true;
    }
}

$db = new DB();
$db->query("SELECT * FROM users");
$db->query("INSERT INTO users (name) VALUES ('test')");
$db->query("UPDATE users SET name = 'updated' WHERE id = 1");
