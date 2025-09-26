<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { 
    currentStepData, 
    visualizationTheme,
    visualizationUtils 
  } from '../../stores/visualization';
  import type { VisualizationStep } from '../../types/visualization';

  // 组件属性
  export let width = 800;
  export let height = 400;
  export let showIndices = true;
  export let showValues = true;
  export let animationDuration = 300;

  // DOM引用
  let canvas: HTMLCanvasElement;
  let ctx: CanvasRenderingContext2D;

  // 本地状态
  let animationFrame: number | null = null;
  let isAnimating = false;

  // 响应式语句
  $: if (ctx && $currentStepData) {
    drawVisualization($currentStepData);
  }

  // 生命周期
  onMount(() => {
    ctx = canvas.getContext('2d')!;
    setupCanvas();
  });

  onDestroy(() => {
    if (animationFrame) {
      cancelAnimationFrame(animationFrame);
    }
  });

  // 设置画布
  function setupCanvas() {
    canvas.width = width;
    canvas.height = height;
    
    // 设置高DPI支持
    const dpr = window.devicePixelRatio || 1;
    const rect = canvas.getBoundingClientRect();
    
    canvas.width = rect.width * dpr;
    canvas.height = rect.height * dpr;
    
    ctx.scale(dpr, dpr);
    canvas.style.width = rect.width + 'px';
    canvas.style.height = rect.height + 'px';
  }

  // 绘制可视化
  function drawVisualization(step: VisualizationStep) {
    if (!ctx || !step) return;

    // 清除画布
    ctx.clearRect(0, 0, width, height);

    // 根据数据类型绘制
    if (Array.isArray(step.data)) {
      drawArrayVisualization(step);
    } else if (isGraphData(step.data)) {
      drawGraphVisualization(step);
    } else if (step.data && typeof step.data === 'object') {
      drawGenericVisualization(step);
    }
  }

  // 绘制数组可视化
  function drawArrayVisualization(step: VisualizationStep) {
    const data = step.data as any[];
    if (!data || data.length === 0) return;

    const padding = 40;
    const availableWidth = width - 2 * padding;
    const availableHeight = height - 2 * padding;
    
    const barWidth = Math.min(availableWidth / data.length, 60);
    const barSpacing = Math.max((availableWidth - barWidth * data.length) / (data.length - 1), 2);
    
    // 计算数值范围
    const maxValue = Math.max(...data.map(v => typeof v === 'number' ? v : 0));
    const minValue = Math.min(...data.map(v => typeof v === 'number' ? v : 0));
    const valueRange = maxValue - minValue || 1;

    // 绘制每个元素
    data.forEach((value, index) => {
      const x = padding + index * (barWidth + barSpacing);
      const normalizedValue = typeof value === 'number' ? (value - minValue) / valueRange : 0.5;
      const barHeight = normalizedValue * availableHeight * 0.7;
      const y = height - padding - barHeight;

      // 确定颜色
      let color = $visualizationTheme.primaryColor;
      
      if (step.highlights && step.highlights.includes(index)) {
        color = $visualizationTheme.highlightColor;
      }
      
      // 检查比较操作
      if (step.comparisons) {
        for (const comparison of step.comparisons) {
          if (comparison.index1 === index || comparison.index2 === index) {
            color = $visualizationTheme.compareColor;
            break;
          }
        }
      }

      // 绘制柱状图
      ctx.fillStyle = color;
      ctx.fillRect(x, y, barWidth, barHeight);

      // 绘制边框
      ctx.strokeStyle = $visualizationTheme.borderColor;
      ctx.lineWidth = 1;
      ctx.strokeRect(x, y, barWidth, barHeight);

      // 绘制索引
      if (showIndices) {
        ctx.fillStyle = $visualizationTheme.textColor;
        ctx.font = '12px Arial';
        ctx.textAlign = 'center';
        ctx.fillText(index.toString(), x + barWidth / 2, height - padding + 15);
      }

      // 绘制数值
      if (showValues) {
        ctx.fillStyle = $visualizationTheme.textColor;
        ctx.font = '14px Arial';
        ctx.textAlign = 'center';
        ctx.fillText(value.toString(), x + barWidth / 2, y - 5);
      }
    });

    // 绘制比较箭头
    if (step.comparisons) {
      step.comparisons.forEach(comparison => {
        drawComparisonArrow(comparison.index1, comparison.index2, barWidth, barSpacing, padding);
      });
    }
  }

  // 绘制比较箭头
  function drawComparisonArrow(index1: number, index2: number, barWidth: number, barSpacing: number, padding: number) {
    const x1 = padding + index1 * (barWidth + barSpacing) + barWidth / 2;
    const x2 = padding + index2 * (barWidth + barSpacing) + barWidth / 2;
    const y = height - padding - 30;

    ctx.strokeStyle = $visualizationTheme.compareColor;
    ctx.lineWidth = 2;
    
    // 绘制连接线
    ctx.beginPath();
    ctx.moveTo(x1, y);
    ctx.lineTo(x2, y);
    ctx.stroke();

    // 绘制箭头
    const arrowSize = 8;
    const angle = Math.atan2(0, x2 - x1);
    
    ctx.beginPath();
    ctx.moveTo(x2, y);
    ctx.lineTo(x2 - arrowSize * Math.cos(angle - Math.PI / 6), y - arrowSize * Math.sin(angle - Math.PI / 6));
    ctx.moveTo(x2, y);
    ctx.lineTo(x2 - arrowSize * Math.cos(angle + Math.PI / 6), y - arrowSize * Math.sin(angle + Math.PI / 6));
    ctx.stroke();
  }

  function isGraphData(data: any): boolean {
    return data && Array.isArray(data.nodes) && Array.isArray(data.edges);
  }

  // 绘制图可视化
  function drawGraphVisualization(step: VisualizationStep) {
    const data = step.data as any;
    const nodes = data.nodes as Array<{ id: string; label: string; x?: number; y?: number; value?: any }>;
    const edges = data.edges as Array<{ from: string; to: string; weight?: any; label?: string }>;

    // 使用力导向布局计算节点位置
    const positions = calculateForceDirectedLayout(nodes, edges);

    // 绘制边
    ctx.strokeStyle = $visualizationTheme.borderColor;
    ctx.lineWidth = 1.5;
    ctx.font = '12px Arial';

    edges.forEach((e, edgeIndex) => {
      const p1 = positions[e.from];
      const p2 = positions[e.to];
      if (!p1 || !p2) return;

      // 检查边是否被高亮
      const highlighted = step.highlights && step.highlights.includes(`edge_${edgeIndex}`);
      ctx.strokeStyle = highlighted ? $visualizationTheme.highlightColor : $visualizationTheme.borderColor;

      // 绘制边
      ctx.beginPath();
      ctx.moveTo(p1.x, p1.y);
      ctx.lineTo(p2.x, p2.y);
      ctx.stroke();

      // 如果是有向图，绘制箭头
      if (data.type === 'directed') {
        drawArrow(p1.x, p1.y, p2.x, p2.y, 16);
      }

      // 显示权重
      if (e.weight !== undefined && e.weight !== null) {
        const midX = (p1.x + p2.x) / 2;
        const midY = (p1.y + p2.y) / 2;

        // 绘制权重背景
        const weightText = e.weight.toString();
        const textMetrics = ctx.measureText(weightText);
        const padding = 4;

        ctx.fillStyle = 'rgba(255, 255, 255, 0.8)';
        ctx.fillRect(
          midX - textMetrics.width / 2 - padding,
          midY - 6 - padding,
          textMetrics.width + padding * 2,
          12 + padding * 2
        );

        // 绘制权重文本
        ctx.fillStyle = $visualizationTheme.textColor;
        ctx.textAlign = 'center';
        ctx.fillText(weightText, midX, midY + 4);
      }
    });

    // 绘制节点
    nodes.forEach((n, i) => {
      const pos = positions[n.id];
      if (!pos) return;

      const radius = 20;
      const highlighted = step.highlights && step.highlights.includes(i);

      // 检查节点是否在比较操作中
      let nodeColor = $visualizationTheme.primaryColor;
      if (highlighted) {
        nodeColor = $visualizationTheme.highlightColor;
      } else if (step.comparisons) {
        for (const comparison of step.comparisons) {
          if (comparison.index1 === i || comparison.index2 === i) {
            nodeColor = $visualizationTheme.compareColor;
            break;
          }
        }
      }

      // 绘制节点圆圈
      ctx.fillStyle = nodeColor;
      ctx.beginPath();
      ctx.arc(pos.x, pos.y, radius, 0, 2 * Math.PI);
      ctx.fill();

      // 绘制节点边框
      ctx.strokeStyle = $visualizationTheme.borderColor;
      ctx.lineWidth = 2;
      ctx.stroke();

      // 绘制节点标签
      ctx.fillStyle = 'white';
      ctx.font = 'bold 12px Arial';
      ctx.textAlign = 'center';
      ctx.fillText(n.label || n.id, pos.x, pos.y + 4);

      // 在节点下方显示值（如果有）
      if (n.value !== undefined && n.value !== null) {
        ctx.fillStyle = $visualizationTheme.textColor;
        ctx.font = '10px Arial';
        ctx.fillText(`(${n.value})`, pos.x, pos.y + radius + 15);
      }
    });
  }

  // 力导向布局算法
  function calculateForceDirectedLayout(nodes: any[], edges: any[]): Record<string, { x: number; y: number }> {
    const positions: Record<string, { x: number; y: number }> = {};
    const velocities: Record<string, { vx: number; vy: number }> = {};

    // 初始化位置（圆形排列）
    const centerX = width / 2;
    const centerY = height / 2;
    const radius = Math.min(width, height) * 0.3;

    nodes.forEach((node, i) => {
      if (typeof node.x === 'number' && typeof node.y === 'number') {
        // 如果节点已有坐标，使用现有坐标
        positions[node.id] = { x: node.x, y: node.y };
      } else {
        // 否则使用圆形布局
        const angle = (2 * Math.PI * i) / nodes.length;
        positions[node.id] = {
          x: centerX + radius * Math.cos(angle),
          y: centerY + radius * Math.sin(angle)
        };
      }
      velocities[node.id] = { vx: 0, vy: 0 };
    });

    // 力导向布局参数
    const k = Math.sqrt((width * height) / nodes.length);
    const iterations = 50;
    const cooling = 0.99;
    let temperature = width / 10;

    // 迭代计算
    for (let iter = 0; iter < iterations; iter++) {
      // 计算排斥力
      nodes.forEach(nodeA => {
        velocities[nodeA.id] = { vx: 0, vy: 0 };

        nodes.forEach(nodeB => {
          if (nodeA.id !== nodeB.id) {
            const dx = positions[nodeA.id].x - positions[nodeB.id].x;
            const dy = positions[nodeA.id].y - positions[nodeB.id].y;
            const distance = Math.sqrt(dx * dx + dy * dy) || 1;

            if (distance < k * 2) {
              const repulsiveForce = k * k / distance;
              velocities[nodeA.id].vx += (dx / distance) * repulsiveForce;
              velocities[nodeA.id].vy += (dy / distance) * repulsiveForce;
            }
          }
        });
      });

      // 计算吸引力（基于边连接）
      edges.forEach(edge => {
        const dx = positions[edge.to].x - positions[edge.from].x;
        const dy = positions[edge.to].y - positions[edge.from].y;
        const distance = Math.sqrt(dx * dx + dy * dy) || 1;

        const attractiveForce = distance * distance / k;
        const fx = (dx / distance) * attractiveForce;
        const fy = (dy / distance) * attractiveForce;

        velocities[edge.from].vx += fx;
        velocities[edge.from].vy += fy;
        velocities[edge.to].vx -= fx;
        velocities[edge.to].vy -= fy;
      });

      // 更新位置
      nodes.forEach(node => {
        const v = velocities[node.id];
        const speed = Math.sqrt(v.vx * v.vx + v.vy * v.vy);

        if (speed > temperature) {
          v.vx = (v.vx / speed) * temperature;
          v.vy = (v.vy / speed) * temperature;
        }

        positions[node.id].x += v.vx;
        positions[node.id].y += v.vy;

        // 边界约束
        const margin = 50;
        positions[node.id].x = Math.max(margin, Math.min(width - margin, positions[node.id].x));
        positions[node.id].y = Math.max(margin, Math.min(height - margin, positions[node.id].y));
      });

      temperature *= cooling;
    }

    return positions;
  }

  // 绘制箭头
  function drawArrow(fromX: number, fromY: number, toX: number, toY: number, nodeRadius: number) {
    const dx = toX - fromX;
    const dy = toY - fromY;
    const distance = Math.sqrt(dx * dx + dy * dy);

    if (distance === 0) return;

    // 计算箭头起点（避开节点圆圈）
    const unitX = dx / distance;
    const unitY = dy / distance;
    const endX = toX - unitX * nodeRadius;
    const endY = toY - unitY * nodeRadius;

    // 箭头参数
    const arrowLength = 12;
    const arrowAngle = Math.PI / 6;

    // 计算箭头两个边的端点
    const angle = Math.atan2(dy, dx);
    const arrowX1 = endX - arrowLength * Math.cos(angle - arrowAngle);
    const arrowY1 = endY - arrowLength * Math.sin(angle - arrowAngle);
    const arrowX2 = endX - arrowLength * Math.cos(angle + arrowAngle);
    const arrowY2 = endY - arrowLength * Math.sin(angle + arrowAngle);

    // 绘制箭头
    ctx.beginPath();
    ctx.moveTo(endX, endY);
    ctx.lineTo(arrowX1, arrowY1);
    ctx.moveTo(endX, endY);
    ctx.lineTo(arrowX2, arrowY2);
    ctx.stroke();
  }

  // 绘制通用可视化
  function drawGenericVisualization(step: VisualizationStep) {
    // 简单的文本显示
    ctx.fillStyle = $visualizationTheme.textColor;
    ctx.font = '16px Arial';
    ctx.textAlign = 'center';
    ctx.fillText('数据可视化', width / 2, height / 2);
    ctx.fillText(JSON.stringify(step.data).substring(0, 50) + '...', width / 2, height / 2 + 30);
  }

  // 动画效果（预留）
  function animateTransition(fromStep: VisualizationStep, toStep: VisualizationStep) {
    if (isAnimating) return;
    
    isAnimating = true;
    const startTime = performance.now();
    
    const animate = (currentTime: number) => {
      const elapsed = currentTime - startTime;
      const progress = Math.min(elapsed / animationDuration, 1);
      
      // 这里可以添加插值动画逻辑
      
      if (progress < 1) {
        animationFrame = requestAnimationFrame(animate);
      } else {
        isAnimating = false;
        drawVisualization(toStep);
      }
    };
    
    animationFrame = requestAnimationFrame(animate);
  }

  // 导出画布为图片
  export function exportAsImage(): string {
    return canvas.toDataURL('image/png');
  }

  // 重置画布大小
  export function resize(newWidth: number, newHeight: number) {
    width = newWidth;
    height = newHeight;
    setupCanvas();
    if ($currentStepData) {
      drawVisualization($currentStepData);
    }
  }
</script>

<div class="visualization-canvas">
  <canvas
    bind:this={canvas}
    {width}
    {height}
    style="width: {width}px; height: {height}px;"
  ></canvas>
  
  {#if $currentStepData}
    <div class="step-info">
      <div class="step-description">
        {$currentStepData.description}
      </div>
      
      {#if visualizationUtils.hasComparisons($currentStepData)}
        <div class="comparisons">
          {#each $currentStepData.comparisons as comparison}
            <span class="comparison">
              比较 [{comparison.index1}] {visualizationUtils.getComparisonResultText(comparison.result)} [{comparison.index2}]
            </span>
          {/each}
        </div>
      {/if}
      
      {#if visualizationUtils.hasOperations($currentStepData)}
        <div class="operations">
          {#each $currentStepData.operations as operation}
            <span class="operation">
              {visualizationUtils.getOperationTypeName(operation.type)}: {operation.description}
            </span>
          {/each}
        </div>
      {/if}
    </div>
  {/if}
</div>

<style>
  .visualization-canvas {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  canvas {
    border: 1px solid #e5e7eb;
    border-radius: 4px;
    background: white;
  }

  .step-info {
    width: 100%;
    max-width: 600px;
    text-align: center;
  }

  .step-description {
    font-size: 1rem;
    color: #1f2937;
    margin-bottom: 0.5rem;
    font-weight: 500;
  }

  .comparisons,
  .operations {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    justify-content: center;
    margin-top: 0.5rem;
  }

  .comparison,
  .operation {
    padding: 0.25rem 0.5rem;
    background: #f3f4f6;
    border-radius: 4px;
    font-size: 0.875rem;
    color: #4b5563;
  }

  .comparison {
    background: #fef2f2;
    color: #dc2626;
  }

  .operation {
    background: #f0f9ff;
    color: #0369a1;
  }
</style>
