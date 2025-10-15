import numpy as np
from dataclasses import dataclass
from typing import List, Tuple
import matplotlib.pyplot as plt

@dataclass
class Point:
    x: float
    y: float

@dataclass
class Wall:
    start: Point
    end: Point
    
    def is_vertical(self) -> bool:
        return abs(self.start.x - self.end.x) < 1e-6

class RayLaunching:
    def __init__(self, room_size: float = 20.0, step: float = 0.1):
        self.room_size = room_size
        self.step = step
        self.grid_size = int(room_size / step)
        self.power_map = np.zeros((self.grid_size, self.grid_size))
        
        self.walls = [
            Wall(Point(0, 0), Point(room_size, 0)),  
            Wall(Point(room_size, 0), Point(room_size, room_size)),  
            Wall(Point(room_size, room_size), Point(0, room_size)),  
            Wall(Point(0, room_size), Point(0, 0)) 
        ]
        
    def calculate_intersection(self, ray_start: Point, ray_angle: float, wall: Wall) -> Tuple[bool, Point]:
        dx = np.cos(ray_angle)
        dy = np.sin(ray_angle)
        
        wall_dx = wall.end.x - wall.start.x
        wall_dy = wall.end.y - wall.start.y
        
        det = dx * (-wall_dy) - dy * (-wall_dx)
        
        if abs(det) < 1e-6: 
            return False, Point(0, 0)
            
        t = ((wall.start.x - ray_start.x) * (-wall_dy) - (wall.start.y - ray_start.y) * (-wall_dx)) / det
        u = (dx * (wall.start.y - ray_start.y) - dy * (wall.start.x - ray_start.x)) / det
        
        if t >= 0 and 0 <= u <= 1:
            intersection_x = ray_start.x + t * dx
            intersection_y = ray_start.y + t * dy
            return True, Point(intersection_x, intersection_y)
            
        return False, Point(0, 0)
    
    def calculate_reflection_angle(self, incident_angle: float, wall: Wall) -> float:
        if wall.is_vertical():
            return np.pi - incident_angle
        return -incident_angle
    
    def simulate(self, source: Point, power: float = 1.0, angle_step: float = np.pi/180):
        self.source = source  
        for angle in np.arange(0, 2*np.pi, angle_step):
            self.trace_ray(source, angle, power)
    
    def trace_ray(self, start: Point, angle: float, power: float):
        current_point = start
        current_angle = angle
        current_power = power
        
        while current_power > 0.01: 
          
            min_distance = float('inf')
            closest_wall = None
            intersection_point = None
            
            for wall in self.walls:
                has_intersection, point = self.calculate_intersection(current_point, current_angle, wall)
                if has_intersection:
                    distance = np.sqrt((point.x - current_point.x)**2 + (point.y - current_point.y)**2)
                    if distance < min_distance:
                        min_distance = distance
                        closest_wall = wall
                        intersection_point = point
            
         
            self.add_power_along_ray(current_point, intersection_point, current_power)
            
            if closest_wall is None:
                break
                
        
            reflection_angle = self.calculate_reflection_angle(current_angle, closest_wall)
        
            current_point = intersection_point
            current_angle = reflection_angle
            current_power *= 0.5  
    
    def add_power_along_ray(self, start: Point, end: Point, power: float):
        if end is None:
            return
            
        distance = np.sqrt((end.x - start.x)**2 + (end.y - start.y)**2)
        steps = int(distance / self.step)
        
        for i in range(steps):
            t = i / steps
            x = start.x + t * (end.x - start.x)
            y = start.y + t * (end.y - start.y)
            
            grid_x = int(x / self.step)
            grid_y = int(y / self.step)
            
            if 0 <= grid_x < self.grid_size and 0 <= grid_y < self.grid_size:
              
                distance_from_start = np.sqrt((x - start.x)**2 + (y - start.y)**2)
                local_power = power / (1 + distance_from_start)  #
                self.power_map[grid_y, grid_x] += local_power

    def visualize(self):
        plt.figure(figsize=(10, 10))
        
        plt.imshow(self.power_map, 
                  extent=[0, self.room_size, 0, self.room_size],
                  origin='lower',
                  cmap='jet')
        plt.colorbar(label='Moc sygnału')
        
        for wall in self.walls:
            plt.plot([wall.start.x, wall.end.x], 
                    [wall.start.y, wall.end.y], 
                    'k-', linewidth=2, label='Ściany')
        
        plt.plot(self.source.x, self.source.y, 'r*', markersize=15, label='Nadajnik')
        
        plt.title('Mapa mocy sygnału')
        plt.xlabel('X [m]')
        plt.ylabel('Y [m]')
        plt.grid(True)
        plt.legend()
        plt.show()

def main():
    simulator = RayLaunching(room_size=20.0, step=0.1)
    
    source = Point(10.0, 10.0)
    
    simulator.simulate(source, power=1.0, angle_step=np.pi/180)
    
    simulator.visualize()

if __name__ == "__main__":
    main()