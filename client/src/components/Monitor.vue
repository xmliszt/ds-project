<template>
  <div class="monitor">
    <div>0</div>
    <div class="minitor-content-wrapper">
      <div class="canvas-left-placeholder">
        <el-table :data="virtualNodeInfoData" height="38rem" size="mini" border>
          <el-table-column prop="name" label="Node Name" width="120" sortable>
          </el-table-column>
          <el-table-column
            prop="location"
            label="Location"
            width="120"
            sortable
          ></el-table-column>
          <el-table-column label="From">
            <template slot-scope="scope">
              Node {{ scope.row.physicalNodeID }}
            </template>
          </el-table-column>
          <el-table-column label="Alive">
            <template slot-scope="scope">
              <el-tag
                type="success"
                effect="plain"
                v-if="scope.row.alive"
                size="mini"
                >Alive</el-tag
              >
              <el-tag type="danger" effect="dark" v-else size="mini"
                >Dead</el-tag
              >
            </template>
          </el-table-column>
        </el-table>
      </div>
      <canvas id="ring" width="600" height="600"></canvas>
      <div class="canvas-right-placeholder">
        <el-table
          :data="monitorNodesData"
          height="38rem"
          border
          size="mini"
          :show-header="false"
        >
          <el-table-column>
            <template slot-scope="scope">
              <el-tag
                size="mini"
                v-if="scope.row.s"
                type="success"
                effect="plain"
                >Physical Node {{ scope.row.v }} is running</el-tag
              >
              <el-tag size="mini" v-else type="danger" effect="dark"
                >Physical Node {{ scope.row.v }} is down!</el-tag
              >
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script>
import Chart from "chart.js/auto";

import { getRandomColor, deepCompareArray } from "../util/util";
import { getMonitorInfo } from "../service/monitor";

export default {
  data() {
    return {
      monitorInterval: null,
      ring: null,
      virtualNodeInfoData: [],
      monitorNodesData: [],
      virtualNodes: [],
      virtualNodesProportions: [],
      virtualNodesColors: [],
      physicalNodesColors: Object(),
    };
  },
  methods: {
    async fetchData() {
      var virtualNodeInfoData = [];
      var monitorNodesData = [];
      var virtualNodes = [];
      var virtualNodesProportions = [];
      var virtualNodesColors = [];

      var data = await getMonitorInfo();
      var heartbeatTable = data.data.HeartbeatTable;
      var nodesCapacity = data.data.VirtualNodesCapacity;
      var nodesLocation = data.data.VirtualNodesLocation;
      var nameLocations = data.data.VirtualNodes;

      Object.entries(heartbeatTable).forEach(([nid, alive]) => {
        monitorNodesData.push({
          v: nid,
          s: alive,
        });
        if (!(nid in this.physicalNodesColors)) {
          this.physicalNodesColors[nid] = getRandomColor();
        }
      });

      nameLocations.forEach((name) => {
        var nid = Number(name.split("-")[0]);
        var capacity = nodesCapacity[name];
        var location = nodesLocation[name];
        var alive = heartbeatTable[nid];

        virtualNodes.push(name);
        virtualNodesProportions.push(capacity);
        virtualNodesColors.push(this.physicalNodesColors[nid]);
        virtualNodeInfoData.push({
          name: name,
          location: location,
          physicalNodeID: nid,
          alive: alive,
        });
      });

      if (!deepCompareArray(this.virtualNodeInfoData, virtualNodeInfoData)) {
        this.virtualNodeInfoData = virtualNodeInfoData;
      }

      if (!deepCompareArray(this.monitorNodesData, monitorNodesData)) {
        this.monitorNodesData = monitorNodesData;
      }

      if (
        !deepCompareArray(this.virtualNodes, virtualNodes) ||
        !deepCompareArray(
          this.virtualNodesProportions,
          virtualNodesProportions
        ) ||
        !deepCompareArray(this.virtualNodesColors, virtualNodesColors)
      ) {
        if (!deepCompareArray(this.virtualNodes, virtualNodes)) {
          this.virtualNodes = virtualNodes;
        }

        if (
          !deepCompareArray(
            this.virtualNodesProportions,
            virtualNodesProportions
          )
        ) {
          this.virtualNodesProportions = virtualNodesProportions;
        }

        if (!deepCompareArray(this.virtualNodesColors, virtualNodesColors)) {
          this.virtualNodesColors = virtualNodesColors;
        }

        if (this.ring !== null) {
          this.ring.destroy();
        }

        var ctx = document.getElementById("ring").getContext("2d");
        this.ring = new Chart(ctx, {
          type: "doughnut",
          data: {
            labels: this.virtualNodes,
            datasets: [
              {
                label: "Virtual Nodes Ring Structure",
                data: this.virtualNodesProportions,
                backgroundColor: this.virtualNodesColors,
                hoverOffset: 4,
              },
            ],
          },
          options: {
            responsive: false,
            plugins: {
              tooltip: {
                callbacks: {
                  label: function (context) {
                    var label = context.label;
                    return "Virtual Node: " + label;
                  },
                },
              },
              legend: {
                display: false,
              },
            },
          },
        });
      }
    },
  },
  mounted() {
    if (this.monitorInterval === null) {
      this.monitorInterval = setInterval(() => {
        this.fetchData();
      }, 1000);
    }
  },
  destroyed() {
    if (this.monitorInterval) {
      this.monitorInterval.stopInterval();
    }
  },
};
</script>

<style scoped>
.monitor {
  margin-top: 4rem;
}

.minitor-content-wrapper {
  display: flex;
  justify-content: space-between;
}

.canvas-left-placeholder {
  margin-right: 4rem;
  width: 26rem;
  padding-left: 2rem;
}

.canvas-right-placeholder {
  margin-left: 4rem;
  display: flex;
  align-items: center;
  justify-content: right;
  width: 26rem;
  padding-right: 2rem;
}
</style>
