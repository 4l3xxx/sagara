import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Work } from './entities/work.entity';

@Injectable()
export class WorksService {
  constructor(@InjectRepository(Work) private repo: Repository<Work>) {}

  async findAll(): Promise<any> {
    const works = await this.repo.find({ where: { published: true }, order: { createdAt: 'DESC' } });
    return { success: true, data: works };
  }

  async findBySlug(slug: string): Promise<any> {
    const work = await this.repo.findOne({ where: { slug, published: true } });
    if (!work) throw new NotFoundException(`Work "${slug}" not found`);
    return { success: true, data: work };
  }
}
