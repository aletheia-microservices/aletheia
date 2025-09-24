import matplotlib.pyplot as plt
import numpy as np

apps = ["Digota", "Sockshop", "PostNotification", "SocialNetwork", "MediaMicroservices", "TrainTicket"]
complexity = np.array([4+5, 7+6, 3+2, 10+10, 5+4, 30+21])

complexity = (complexity - complexity.min()) / (complexity.max() - complexity.min())

exec_times = {
    "C1 (Ref. Integrity, Absence of Cascade Delete)": [10, 12, 8, 15, 11, 40],
    "C2 (Ref. Integrity, Concurrent Associations)": [11, 13, 9, 16, 12, 42],
    "C3 (Ref. Integrity, Uncoordinated Replication)": [14, 15, 11, 20, 14, 48],
    "C4 (Ent. Integrity, Uncoordinated Replication)": [9, 11, 7, 14, 10, 39],
    "C5 (Unicity, Conflict Resolution)": [8, 10, 6, 12, 9, 35],
}

plt.figure(figsize=(8,6))

pastel_blue   = "#6baed6"
pastel_green  = "#74c476"
pastel_red    = "#fb6a4a"

plt.plot(complexity, exec_times["C1 (Ref. Integrity, Absence of Cascade Delete)"], 
         color=pastel_blue, linestyle="-", marker="o", label="C1 (Ref. Integrity, Absence of Cascade Delete)")
plt.plot(complexity, exec_times["C2 (Ref. Integrity, Concurrent Associations)"], 
         color=pastel_blue, linestyle="--", marker="s", label="C2 (Ref. Integrity, Concurrent Associations)")
plt.plot(complexity, exec_times["C3 (Ref. Integrity, Uncoordinated Replication)"], 
         color=pastel_blue, linestyle=":", marker="^", label="C3 (Ref. Integrity, Uncoordinated Replication)")

plt.plot(complexity, exec_times["C4 (Ent. Integrity, Uncoordinated Replication)"], 
         color=pastel_green, linestyle="-", marker="D", label="C4 (Ent. Integrity, Uncoordinated Replication)")

plt.plot(complexity, exec_times["C5 (Unicity, Conflict Resolution)"], 
         color=pastel_red, linestyle="-", marker="v", label="C5 (Unicity, Conflict Resolution)")

plt.title("Analysis Time per Application Complexity for each Class of Executions")
plt.xlabel("Normalized Application Complexity")
plt.ylabel("Analysis Time (s)")
plt.legend(title="Classes of Execution")
plt.grid(True, linestyle="--", alpha=0.6)
plt.tight_layout()
plt.show()
