package com.example.kafkagluepoc.kafka;

import com.example.kafkagluepoc.UserSignedUp;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

import java.time.Instant;
import java.util.UUID;

@Component
public class ProducerSimulator {
    private final KafkaProducerService producer;

    public ProducerSimulator(KafkaProducerService producer) {
        this.producer = producer;
    }

    @Scheduled(fixedRate = 2000)
    public void sendEvent() {
        UserSignedUp user = new UserSignedUp();
        user.setUserId(UUID.randomUUID().toString());
        if (Math.random() < 0.5) {
            user.setEmail(null);
        } else {
            user.setEmail("test@example.com");
        }
        user.setTimestamp(Instant.now().toString());
        producer.sendUserSignedUpEvent(user);
    }
}
