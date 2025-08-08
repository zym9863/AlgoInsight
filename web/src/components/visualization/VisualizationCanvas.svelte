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
    } else if (step.data && typeof step.data === 'object') {
      // 处理其他数据类型（图、树等）
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
