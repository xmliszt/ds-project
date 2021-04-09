<template>
  <div class="monitor">
    <div class="minitor-content-wrapper">
      <div class="canvas-left-placeholder">
        <el-table
          :data="virtualNodeInfoData"
          height="38rem"
          v-loading="loading"
          size="mini"
          border
        >
          <el-table-column prop="name" label="Node Name" width="120">
          </el-table-column>
          <el-table-column prop="location" label="Location"></el-table-column>
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
          v-loading="loading"
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

import { getRandomColor } from "../util/util";

export default {
  data() {
    return {
      ring: null,
      loading: false,
      virtualNodeInfoData: [
        {
          name: "1-1",
          location: 100,
          physicalNodeID: 1,
          alive: true,
        },
      ],
      monitorNodesData: [
        {
          v: "1",
          s: false,
        },
        {
          v: "2",
          s: true,
        },
        {
          v: "3",
          s: true,
        },
      ],
      virtualNodes: [
        "1-1",
        "2-1",
        "3-1",
        "1-2",
        "1-3",
        "2-2",
        "3-3",
        "2-3",
        "3-2",
      ],
      virtualNodesProportions: [300, 50, 100, 12, 45, 60, 131, 214, 495],
      virtualNodesColors: [],
    };
  },
  methods: {
    async fetchVirtualNodesInformation() {},
  },
  async mounted() {
    // get nodes information
    this.loading = true;
    await this.fetchVirtualNodesInformation();
    this.loading = false;

    // generate random colors
    for (let i = 0; i < this.virtualNodes.length; i++) {
      this.virtualNodesColors.push(getRandomColor());
    }

    // generate doughnut chart
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
        },
      },
    });
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
  width: 25rem;
  padding-left: 2rem;
}

.canvas-right-placeholder {
  margin-left: 4rem;
  display: flex;
  align-items: center;
  justify-content: right;
  width: 25rem;
  padding-right: 2rem;
}
</style>
