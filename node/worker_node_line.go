package node

import (
	"go.uber.org/zap"
	"teamide/pkg/util"
)

func (this_ *Worker) appendNodeLineList(loadedIdList *[]string, lineList *[][]string, parentLine []string, nodeList []*Info, nodeIdConnNodeIdListCache map[string][]string) {

	for _, one := range nodeList {
		var line []string
		line = append(line, parentLine...)

		if util.ContainsString(line, one.Id) >= 0 {
			continue
		}
		line = append(line, one.Id)

		*lineList = append(*lineList, line)

		if util.ContainsString(*loadedIdList, one.Id) >= 0 {
			continue
		}
		*loadedIdList = append(*loadedIdList, one.Id)

		var connNodeIdList = nodeIdConnNodeIdListCache[one.Id]
		var children = this_.findNodeList(connNodeIdList)
		this_.appendNodeLineList(loadedIdList, lineList, line, children, nodeIdConnNodeIdListCache)

		var parentList []*Info
		for cacheNodeId, cacheConnNodeIdList := range nodeIdConnNodeIdListCache {
			if util.ContainsString(cacheConnNodeIdList, one.Id) >= 0 {
				var cacheNode = this_.findNode(cacheNodeId)
				if cacheNode != nil {
					parentList = append(parentList, cacheNode)
				}
			}
		}
		this_.appendNodeLineList(loadedIdList, lineList, line, parentList, nodeIdConnNodeIdListCache)
	}
}

/**




 */
func (this_ *Worker) findNodeLineList(nodeId string, nodeIdConnNodeIdListCache map[string][]string) (lineList [][]string) {
	Logger.Info("查询节点所有线", zap.Any("nodeId", nodeId))

	var loadedIdList []string
	loadedIdList = append(loadedIdList, nodeId)
	var line []string
	line = append(line, nodeId)

	var find = this_.findNode(nodeId)
	if find != nil {
		var connNodeIdList = nodeIdConnNodeIdListCache[find.Id]
		var children = this_.findNodeList(connNodeIdList)
		this_.appendNodeLineList(&loadedIdList, &lineList, line, children, nodeIdConnNodeIdListCache)
	}
	var parentList []*Info
	for cacheNodeId, cacheConnNodeIdList := range nodeIdConnNodeIdListCache {
		if util.ContainsString(cacheConnNodeIdList, nodeId) >= 0 {
			var cacheNode = this_.findNode(cacheNodeId)
			if cacheNode != nil {
				parentList = append(parentList, cacheNode)
			}
		}
	}
	this_.appendNodeLineList(&loadedIdList, &lineList, line, parentList, nodeIdConnNodeIdListCache)

	return
}

func (this_ *Worker) getNodeLineByFromTo(fromNodeId, toNodeId string, nodeIdConnNodeIdListCache map[string][]string) (lineIdList []string) {

	Logger.Info("查询节点线", zap.Any("fromNodeId", fromNodeId), zap.Any("toNodeId", toNodeId))

	if fromNodeId == toNodeId {
		lineIdList = append(lineIdList, fromNodeId)
		return
	}

	var lineList = this_.findNodeLineList(fromNodeId, nodeIdConnNodeIdListCache)

	for _, line := range lineList {
		Logger.Info("已查询的连线", zap.Any("fromNodeId", fromNodeId), zap.Any("line", line))
		if util.ContainsString(line, toNodeId) >= 0 {
			if len(lineIdList) == 0 || len(line) < len(lineIdList) {
				lineIdList = line
			}
		}
	}
	return
}
