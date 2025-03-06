CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS maps (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,  
    width INT NOT NULL,        
    height INT NOT NULL,       
    bottom_left GEOMETRY(Point, 3857) NOT NULL, 
    top_right GEOMETRY(Point, 3857) NOT NULL
);

CREATE TABLE IF NOT EXISTS walls (
    id SERIAL PRIMARY KEY,
    map_id INT REFERENCES maps(id) ON DELETE CASCADE, 
    bottom_left GEOMETRY(PointZ, 3857) NOT NULL, 
    top_right GEOMETRY(PointZ, 3857) NOT NULL,   
    wall_geometry GEOMETRY(PolygonZ, 3857) NOT NULL
);

