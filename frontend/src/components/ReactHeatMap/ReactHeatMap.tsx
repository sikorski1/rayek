import { OrbitControls } from "@react-three/drei";
import { Canvas } from "@react-three/fiber";

const generateData = () => {
  let data = [];
  for (let x = 0; x < 25; x++) {
    for (let y = 0; y < 25; y++) {
      for (let z = 0; z < 25; z++) {
        data.push({ 
          x:x/5, 
          y:y/5, 
          z:z/5, 
          value: Math.floor(Math.random() * 16) 
        });
      }
    }
  }
  return data;
};
  
  const data = generateData();
export default function ReactHeatMap() {
	return (
		<div style={{height:"100vh", backgroundColor:"white"}}>
			<Canvas camera={{ position: [5, 5, 5] }}>
				<OrbitControls />
				<ambientLight />
				<pointLight position={[10, 10, 10]} />
				{data.map((point, i) => (
					<mesh key={i} position={[point.x, point.y, point.z]}>
						<sphereGeometry args={[0.1, 16, 16]} />
						<meshStandardMaterial color={`hsl(${point.value * 10}, 100%, 50%)`} />
					</mesh>
				))}
			</Canvas>
		</div>
	);
}
