import java.util.List;

public class Worker implements Runnable {

    private final int workerId;
    private final TaskQueue taskQueue;
    private final List<String> results;

    public Worker(int workerId, TaskQueue taskQueue, List<String> results) {
        this.workerId = workerId;
        this.taskQueue = taskQueue;
        this.results = results;
    }
    @Override
    public void run() {
        System.out.println("Worker " + workerId + " started.");
        while (true) {
            Integer task = taskQueue.getTask();
            if (task == null) {
                break;
            }
            try {
                System.out.println("Worker " + workerId + " processing Task " + task);
                Thread.sleep(500);
                int result = task * 2;
                synchronized (results) {
                    results.add("Task " + task + " -> Result " + result);
                }
                System.out.println("Worker " + workerId + " completed Task " + task);

            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                System.err.println("Worker " + workerId + " interrupted: " + e.getMessage());
                break;
            }
        }
        System.out.println("Worker " + workerId + " finished.");
    }
}
