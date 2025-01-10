import os
import sys
import fileinput
import codecs

def replace_in_file(file_path, old_module, new_module):
    # 使用 utf-8 编码读取和写入文件
    with codecs.open(file_path, 'r', encoding='utf-8') as file:
        content = file.read()
    
    content = content.replace(old_module, new_module)
    
    with codecs.open(file_path, 'w', encoding='utf-8') as file:
        file.write(content)

def replace_module_name(old_module, new_module):
    # 替换go.mod文件中的模块名
    if os.path.exists('go.mod'):
        replace_in_file('go.mod', old_module, new_module)
        print(f"已更新 go.mod 文件")
    else:
        print("警告: 未找到 go.mod 文件")
    
    # 遍历所有.go文件并替换导入路径
    files_updated = 0
    for root, _, files in os.walk('.'):
        for file in files:
            if file.endswith('.go'):
                file_path = os.path.join(root, file)
                try:
                    replace_in_file(file_path, old_module, new_module)
                    files_updated += 1
                except Exception as e:
                    print(f"处理文件 {file_path} 时出错: {str(e)}")
    print(f"已更新 {files_updated} 个 .go 文件")

def main():
    if len(sys.argv) != 2:
        print("使用方法: python script.py <新模块名>")
        print("例如: python script.py github.com/YXRRXY/autoAuth")
        sys.exit(1)

    old_module = "github.com/YXRRXY/autoAuth"
    new_module = sys.argv[1]
    
    try:
        print(f"开始替换模块名...")
        print(f"从: {old_module}")
        print(f"到: {new_module}")
        replace_module_name(old_module, new_module)
        print("替换完成!")
    except Exception as e:
        print(f"发生错误: {str(e)}")
        sys.exit(1)

if __name__ == "__main__":
    main() 