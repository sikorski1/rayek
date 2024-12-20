import { useParams } from "react-router-dom"

export default function SingleMap() {
    const {id} = useParams()
    return <div>{id}</div>
}