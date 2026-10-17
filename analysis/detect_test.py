#!/usr/bin/env python3
"""
测试脚本 - 用于诊断 TFLite 和 ONNX 模型的检测问题
使用方法: python detect_test.py <image_path> [model_path]
"""

import sys
import os
import cv2
import numpy as np

# 添加当前目录到路径
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from detect import ObjectDetector, TFLiteBackend, ONNXBackend


def draw_detections(image: np.ndarray, detections: list) -> np.ndarray:
    """在图像上绘制检测结果"""
    img = image.copy()
    for det in detections:
        box = det["box"]
        label = det["label"]
        conf = det["confidence"]

        x1, y1 = box["x_min"], box["y_min"]
        x2, y2 = box["x_max"], box["y_max"]

        # 绘制边框
        cv2.rectangle(img, (x1, y1), (x2, y2), (0, 255, 0), 2)

        # 绘制标签
        text = f"{label}: {conf:.2%}"
        (tw, th), _ = cv2.getTextSize(text, cv2.FONT_HERSHEY_SIMPLEX, 0.6, 1)
        cv2.rectangle(img, (x1, y1 - th - 10), (x1 + tw, y1), (0, 255, 0), -1)
        cv2.putText(
            img, text, (x1, y1 - 5), cv2.FONT_HERSHEY_SIMPLEX, 0.6, (0, 0, 0), 1
        )

    return img


def test_model(model_path: str, image_path: str):
    """测试单个模型"""
    print(f"\n{'='*60}")
    print(f"测试模型: {model_path}")
    print(f"测试图片: {image_path}")
    print(f"{'='*60}")

    # 读取图像
    image = cv2.imread(image_path)
    if image is None:
        print(f"错误: 无法读取图像 {image_path}")
        return

    print(f"图像尺寸: {image.shape}")

    # 创建检测器
    detector = ObjectDetector(model_path)
    if not detector.load_model():
        print("错误: 模型加载失败")
        return

    print(f"模型类型: {detector.model_type}")
    print(f"输入形状: {detector.input_shape}")

    # 检查后端类型
    if isinstance(detector.backend, TFLiteBackend):
        print(f"TFLite 格式: {'NHWC' if detector.backend.is_nhwc() else 'NCHW'}")

        # 调试: 查看输入输出详情
        print(f"输入详情: {detector.backend.input_details}")
        print(f"输出详情: {detector.backend.output_details}")

    # 执行检测
    print("\n执行检测...")

    # 调试: 手动执行预处理并检查
    input_tensor = detector._preprocess(image)
    print(f"预处理后张量形状: {input_tensor.shape}")
    print(f"预处理后张量范围: [{input_tensor.min():.4f}, {input_tensor.max():.4f}]")
    print(f"预处理后张量类型: {input_tensor.dtype}")

    # 执行推理
    if detector.backend is None:
        print("错误: 后端未初始化")
        return
    output = detector.backend.infer(input_tensor)
    print(f"原始输出形状: {output.shape}")
    print(f"原始输出范围: [{output.min():.4f}, {output.max():.4f}]")

    # 检查输出格式
    if len(output.shape) == 3:
        print(
            f"输出维度分析: batch={output.shape[0]}, dim1={output.shape[1]}, dim2={output.shape[2]}"
        )
        # 如果是 (1, 8400, 84) 格式，不需要转置
        # 如果是 (1, 84, 8400) 格式，需要转置
        if output.shape[1] == 84 and output.shape[2] == 8400:
            print("输出格式: (1, 84, 8400) - YOLO 标准格式")
        elif output.shape[1] == 8400 and output.shape[2] == 84:
            print("输出格式: (1, 8400, 84) - 需要调整后处理逻辑")

    # 使用标准检测流程
    detections, inference_time = detector.detect(image, threshold=0.25)

    print(f"\n检测结果 (阈值=0.25):")
    print(f"推理耗时: {inference_time:.2f} ms")
    print(f"检测到 {len(detections)} 个目标:")

    for i, det in enumerate(detections):
        print(f"  {i+1}. {det['label']}: {det['confidence']:.2%} at {det['box']}")
        # 检查置信度是否异常
        if det["confidence"] > 1.0:
            print(f"     ⚠️ 警告: 置信度超过100%! 原始值: {det['confidence']}")

    # 保存结果图像
    model_name = os.path.splitext(os.path.basename(model_path))[0]
    output_path = f"/Users/xugo/Desktop/gowvp/gb28181/analysis/result_{model_name}.jpg"
    result_img = draw_detections(image, detections)
    cv2.imwrite(output_path, result_img)
    print(f"\n结果图像已保存: {output_path}")


def main():
    # 默认图像路径
    image_path = "/Users/xugo/Desktop/gowvp/gb28181/out.png"

    # 命令行参数
    if len(sys.argv) > 1:
        image_path = sys.argv[1]

    # 模型路径
    tflite_model = "/Users/xugo/Desktop/gowvp/gb28181/configs/owl.tflite"
    onnx_model = "/Users/xugo/Desktop/gowvp/gb28181/analysis/owl.onnx"

    if len(sys.argv) > 2:
        # 只测试指定模型
        test_model(sys.argv[2], image_path)
    else:
        # 测试所有可用模型
        if os.path.exists(onnx_model):
            test_model(onnx_model, image_path)

        if os.path.exists(tflite_model):
            test_model(tflite_model, image_path)


if __name__ == "__main__":
    main()
