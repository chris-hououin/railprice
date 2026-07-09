package client

import (
	"fmt"
	"strings"
)

type Cluster struct {
	ClusterId  string
	ClusterNlc string
	end        string
	start      string
}

func newCluster(line string) Cluster {
	clusterId := line[1:5]
	clusterNlc := line[5:9]
	end := line[9:17]
	start := line[17:25]
	c := Cluster{
		clusterId,
		clusterNlc,
		end,
		start,
	}
	return c
}

func determineClusterRecordType(line string) RecordType {
	if strings.HasPrefix(line, "/!!") {
		return CommentRecord
	} else {
		return ClusterRecord
	}
}

func parseClusterLine(line string) interface{} {
	recordType := determineClusterRecordType(line)
	if recordType == ClusterRecord {
		return newCluster(line)
	} else if recordType == CommentRecord {
		return nil
	}
	return nil
}

func parseClusterLines(lines []string) []Cluster {
	var clusters []Cluster
	for l := range lines {
		parsed := parseClusterLine(lines[l])
		switch parsed.(type) {
		case nil:
			break
		case Cluster:
			c := parsed.(Cluster)
			clusters = append(clusters, c)
		}
	}
	fmt.Println("clusters", len(clusters))

	return clusters
}

func ReadClusters(filename string) []Cluster {
	clusterLines, _ := readFileLines(filename + "/" + filename + ".FSC")
	clusters := parseClusterLines(clusterLines)
	return clusters
}

func MapClusters(clusters []Cluster) map[string][]string {
	mappedClusters := make(map[string][]string)
	for _, cluster := range clusters {
		_, exists := mappedClusters[cluster.ClusterId]
		if !exists {
			mappedClusters[cluster.ClusterId] = []string{cluster.ClusterNlc}
		} else {
			mappedClusters[cluster.ClusterId] = append(mappedClusters[cluster.ClusterId], cluster.ClusterNlc)
		}
	}
	return mappedClusters
}

func FetchClusterStations(nlc string, clusterNlcMap map[string][]string) []string {
	_, exists := clusterNlcMap[nlc]
	if exists {
		return clusterNlcMap[nlc]
	} else {
		// The NLC is not a cluster, return it as only station / station-group
		return []string{nlc}
	}

}
