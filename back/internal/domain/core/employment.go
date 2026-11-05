package core

// Employment represents a contract between a employer artifact (Industrial | PublicSystems)
// and a worker artifact (Residential).
type Employment struct {
	EmployerAnchor TilePosition
	WorkerAnchor   TilePosition
	WorkerArtifact ArtifactType
	JobsCount      uint8
}
