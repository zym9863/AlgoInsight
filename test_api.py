#!/usr/bin/env python3
"""
算法洞察平台API测试脚本
用于测试后端API的功能和性能
"""

import requests
import json
import time
import sys

# API基础URL
BASE_URL = "http://localhost:8080/api"

def test_health():
    """测试健康检查接口"""
    print("🔍 测试健康检查接口...")
    try:
        response = requests.get(f"{BASE_URL}/health", timeout=5)
        if response.status_code == 200:
            data = response.json()
            print(f"✅ 健康检查通过: {data['message']}")
            return True
        else:
            print(f"❌ 健康检查失败: HTTP {response.status_code}")
            return False
    except Exception as e:
        print(f"❌ 健康检查异常: {e}")
        return False

def test_get_algorithms():
    """测试获取算法列表"""
    print("\n🔍 测试获取算法列表...")
    try:
        response = requests.get(f"{BASE_URL}/algorithms", timeout=10)
        if response.status_code == 200:
            data = response.json()
            if data.get('success'):
                algorithms = data.get('data', [])
                print(f"✅ 成功获取 {len(algorithms)} 个算法")
                for alg in algorithms[:3]:  # 显示前3个算法
                    print(f"   - {alg['name']} ({alg['category']})")
                return algorithms
            else:
                print(f"❌ API返回失败: {data}")
                return []
        else:
            print(f"❌ 获取算法列表失败: HTTP {response.status_code}")
            return []
    except Exception as e:
        print(f"❌ 获取算法列表异常: {e}")
        return []

def test_get_algorithm_by_category():
    """测试按类别获取算法"""
    print("\n🔍 测试按类别获取算法...")
    categories = ["sorting", "searching", "graph"]
    
    for category in categories:
        try:
            response = requests.get(f"{BASE_URL}/algorithms/category/{category}", timeout=10)
            if response.status_code == 200:
                data = response.json()
                if data.get('success'):
                    algorithms = data.get('data', [])
                    print(f"✅ {category} 类别: {len(algorithms)} 个算法")
                else:
                    print(f"❌ {category} 类别获取失败: {data}")
            else:
                print(f"❌ {category} 类别获取失败: HTTP {response.status_code}")
        except Exception as e:
            print(f"❌ {category} 类别获取异常: {e}")

def test_algorithm_info():
    """测试获取算法详细信息"""
    print("\n🔍 测试获取算法详细信息...")
    algorithm_ids = ["bubble_sort", "quick_sort", "binary_search"]
    
    for alg_id in algorithm_ids:
        try:
            response = requests.get(f"{BASE_URL}/algorithms/info/{alg_id}", timeout=10)
            if response.status_code == 200:
                data = response.json()
                if data.get('success'):
                    alg_info = data.get('data', {})
                    print(f"✅ {alg_info.get('name', alg_id)}: {alg_info.get('timeComplexity', 'N/A')}")
                else:
                    print(f"❌ {alg_id} 信息获取失败: {data}")
            elif response.status_code == 404:
                print(f"⚠️  {alg_id} 算法不存在")
            else:
                print(f"❌ {alg_id} 信息获取失败: HTTP {response.status_code}")
        except Exception as e:
            print(f"❌ {alg_id} 信息获取异常: {e}")

def test_visualization():
    """测试可视化接口"""
    print("\n🔍 测试可视化接口...")
    test_data = {
        "algorithmId": "bubble_sort",
        "data": [64, 34, 25, 12, 22, 11, 90, 5],
        "parameters": {}
    }
    
    try:
        response = requests.post(
            f"{BASE_URL}/visualize/execute", 
            json=test_data, 
            timeout=30
        )
        if response.status_code == 200:
            data = response.json()
            if data.get('success'):
                result = data.get('data', {})
                steps = result.get('totalSteps', 0)
                exec_time = result.get('executionTime', 0)
                print(f"✅ 可视化执行成功: {steps} 步骤, {exec_time/1000000:.2f}ms")
                return True
            else:
                print(f"❌ 可视化执行失败: {data}")
                return False
        else:
            print(f"❌ 可视化请求失败: HTTP {response.status_code}")
            return False
    except Exception as e:
        print(f"❌ 可视化请求异常: {e}")
        return False

def test_data_generation():
    """测试数据生成接口"""
    print("\n🔍 测试数据生成接口...")
    test_request = {
        "dataType": "array",
        "size": 10,
        "pattern": "random"
    }
    
    try:
        response = requests.post(
            f"{BASE_URL}/data/generate", 
            json=test_request, 
            timeout=10
        )
        if response.status_code == 200:
            data = response.json()
            if data.get('success'):
                generated_data = data.get('data', {})
                content = generated_data.get('content', [])
                print(f"✅ 数据生成成功: {len(content)} 个元素")
                print(f"   数据预览: {content[:5]}...")
                return True
            else:
                print(f"❌ 数据生成失败: {data}")
                return False
        else:
            print(f"❌ 数据生成请求失败: HTTP {response.status_code}")
            return False
    except Exception as e:
        print(f"❌ 数据生成请求异常: {e}")
        return False

def test_data_presets():
    """测试预设数据接口"""
    print("\n🔍 测试预设数据接口...")
    try:
        response = requests.get(f"{BASE_URL}/data/presets", timeout=10)
        if response.status_code == 200:
            data = response.json()
            if data.get('success'):
                presets = data.get('data', [])
                print(f"✅ 获取预设数据成功: {len(presets)} 个预设")
                for preset in presets[:2]:  # 显示前2个预设
                    print(f"   - {preset['name']}: {preset['size']} 个元素")
                return True
            else:
                print(f"❌ 获取预设数据失败: {data}")
                return False
        else:
            print(f"❌ 预设数据请求失败: HTTP {response.status_code}")
            return False
    except Exception as e:
        print(f"❌ 预设数据请求异常: {e}")
        return False

def performance_test():
    """简单的性能测试"""
    print("\n🚀 执行性能测试...")
    
    # 测试多次API调用的响应时间
    test_cases = [
        ("健康检查", lambda: requests.get(f"{BASE_URL}/health", timeout=5)),
        ("获取算法列表", lambda: requests.get(f"{BASE_URL}/algorithms", timeout=10)),
        ("获取排序算法", lambda: requests.get(f"{BASE_URL}/algorithms/category/sorting", timeout=10)),
    ]
    
    for name, test_func in test_cases:
        times = []
        for i in range(5):
            start_time = time.time()
            try:
                response = test_func()
                if response.status_code == 200:
                    end_time = time.time()
                    times.append((end_time - start_time) * 1000)  # 转换为毫秒
                else:
                    print(f"❌ {name} 测试失败: HTTP {response.status_code}")
                    break
            except Exception as e:
                print(f"❌ {name} 测试异常: {e}")
                break
        
        if times:
            avg_time = sum(times) / len(times)
            min_time = min(times)
            max_time = max(times)
            print(f"📊 {name}: 平均 {avg_time:.2f}ms, 最小 {min_time:.2f}ms, 最大 {max_time:.2f}ms")

def main():
    """主测试函数"""
    print("🎯 算法洞察平台API测试开始")
    print("=" * 50)
    
    # 基础功能测试
    if not test_health():
        print("❌ 服务器未运行或不可访问，请检查后端服务")
        sys.exit(1)
    
    test_get_algorithms()
    test_get_algorithm_by_category()
    test_algorithm_info()
    test_visualization()
    test_data_generation()
    test_data_presets()
    
    # 性能测试
    performance_test()
    
    print("\n" + "=" * 50)
    print("✅ API测试完成")

if __name__ == "__main__":
    main()
