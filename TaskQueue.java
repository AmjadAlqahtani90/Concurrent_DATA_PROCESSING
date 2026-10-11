import java.util.LinkedList;
import java.util.Queue;

public class TaskQueue {

    private final Queue<Integer> tasks = new LinkedList<>();
    public synchronized void addTask(int task) {
        tasks.add(task);
    }
    public synchronized Integer getTask() {
        return tasks.poll();
    }
}