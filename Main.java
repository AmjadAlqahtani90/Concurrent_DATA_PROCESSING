import java.util.ArrayList;
import java.util.List;
import java.io.BufferedWriter;
import java.io.FileWriter;
import java.io.IOException;

public class Main {
    public static void main(String[] args) {

        TaskQueue taskQueue = new TaskQueue(); // shared resouces 
        List<String> results = new ArrayList<>(); // shared output resources 

        for (int i = 1; i <= 10; i++) {
            taskQueue.addTask(i);
        }
        // worker threads
        Thread t1 = new Thread(new Worker(1, taskQueue, results));
        Thread t2 = new Thread(new Worker(2, taskQueue, results));
        Thread t3 = new Thread(new Worker(3, taskQueue, results));
        
        t1.start();
        t2.start();
        t3.start();
        
        try {
            t1.join();
            t2.join();
            t3.join();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            System.err.println("Main thread interrupted.");
            return;
        }
      
        System.out.println("\nFinal Results:");

        synchronized (results) {
            for (String result : results) {
                System.out.println(result);
            }
        }
        
        try (BufferedWriter writer =
                new BufferedWriter(new FileWriter("results.txt"))) {

            synchronized (results) {
                for (String result : results) {
                    writer.write(result);
                    writer.newLine();
                }
            }
            System.out.println("\nResults saved to results.txt");

        } catch (IOException e) {
            System.err.println("Error writing results: " + e.getMessage());
        }
        System.out.println("Data Processing System finished.");
    }
}