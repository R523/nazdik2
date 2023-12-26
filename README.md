# nazdik -- second generation

## Introduction

Our system is engineered to accurately measure distances using an ultrasonic sensor affixed to the board.
The key to its functionality lies in measuring the time it takes for an ultrasonic pulse to travel to an object and back to the sensor,
thereby calculating the distance based on the speed of sound. Once the distance information is obtained, it is then conveyed in real-time to a React web application via WebSocket.
This full-duplex communication protocol ensures that the web application receives instant updates, allowing for live monitoring and responsive interactions based on the sensed distances
