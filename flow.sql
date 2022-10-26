/*
 Navicat Premium Data Transfer

 Target Server Type    : MySQL
 Target Server Version : 80030
 File Encoding         : 65001

 Date: 25/10/2022 17:54:26
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for QiniuCdnsFlow
-- ----------------------------
DROP TABLE IF EXISTS `QiniuCdnsFlow`;
CREATE TABLE `QiniuCdnsFlow`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `domain` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `date` date NOT NULL,
  `bandwidthmax` int NULL DEFAULT NULL,
  `bytesum` bigint NULL DEFAULT NULL,
  PRIMARY KEY (`id`, `domain`, `date`) USING BTREE,
  INDEX `index`(`domain` ASC, `date` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8410 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for QiniuHubsFlow
-- ----------------------------
DROP TABLE IF EXISTS `QiniuHubsFlow`;
CREATE TABLE `QiniuHubsFlow`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `hub` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `date` date NOT NULL,
  `up` json NULL,
  `down` json NULL,
  `updown` json NULL,
  PRIMARY KEY (`id`, `hub`, `date`) USING BTREE,
  INDEX `QiniuHubsFlow`(`hub` ASC, `date` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7031 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for WangsuCdnFlow
-- ----------------------------
DROP TABLE IF EXISTS `WangsuCdnFlow`;
CREATE TABLE `WangsuCdnFlow`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `channel` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `date` date NOT NULL,
  `peakValue` int NULL DEFAULT NULL,
  `totalFlow` int NULL DEFAULT NULL,
  PRIMARY KEY (`id`, `channel`, `date`) USING BTREE,
  INDEX `channel`(`channel` ASC, `date` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 697 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for WangsuLiveFlow
-- ----------------------------
DROP TABLE IF EXISTS `WangsuLiveFlow`;
CREATE TABLE `WangsuLiveFlow`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `channel` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `date` date NOT NULL,
  `peakValue` int NULL DEFAULT NULL,
  `totalFlow` int NULL DEFAULT NULL,
  PRIMARY KEY (`id`, `channel`, `date`) USING BTREE,
  INDEX `channel`(`channel` ASC, `date` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 363 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
