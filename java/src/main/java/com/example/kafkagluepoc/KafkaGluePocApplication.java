package com.example.kafkagluepoc;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.scheduling.annotation.EnableScheduling;

@SpringBootApplication
@EnableScheduling
public class KafkaGluePocApplication {
    public static void main(String[] args) {
        SpringApplication.run(KafkaGluePocApplication.class, args);
    }
}
